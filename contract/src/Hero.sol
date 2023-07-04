// SPDX-License-Identifier: MIT

/***   
   ________          _          __  __                         
  / ____/ /_  ____ _(_)___     / / / /__  _________  ___  _____
 / /   / __ \/ __ `/ / __ \   / /_/ / _ \/ ___/ __ \/ _ \/ ___/
/ /___/ / / / /_/ / / / / /  / __  /  __/ /  / /_/ /  __(__  ) 
\____/_/ /_/\__,_/_/_/ /_/  /_/ /_/\___/_/   \____/\___/____/  
*/

/// @title Chain Heroes
/// @author rayne (anthonyoliai)

pragma solidity >=0.8.0;

import "./ERC721.sol";
import "./Owned.sol";
import "./MerkleProof.sol";
import "./Strings.sol";

/// @notice Thrown when attempting to mint while total supply has been minted.
error MintedOut();
/// @notice Thrown when minter does not have enough ether.
error NotEnoughFunds();
/// @notice Thrown when a public minter / whitelist minter has reached their mint capacity.
error AlreadyClaimed();
/// @notice Thrown when the hero-sale is not active.
error HeroSaleNotActive();
/// @notice Thrown when the public sale is not active.
error PublicSaleNotActive();
/// @notice Thrown when the msg.sender is not in the hero list.
error NotHeroListed();
/// @notice Thrown when a signer is not authorized.
error NotAuthorized();

contract Hero is ERC721, Owned {
    using Strings for uint256;

    /// @notice The total supply of heroes.
    uint256 public constant MAX_SUPPLY = 8337;
    /// @notice Mint price.
    uint256 public mintPrice = 0.05 ether;
    /// @notice The current supply starts at 10 due to the team minting from tokenID 0 to 14, and 15-24 reserved for legendaries.
    uint256 public totalSupply = 10;

    /// @notice The base URI.
    string baseURI;

    /// @notice Returns true when the hero sale is active, false otherwise.
    bool public heroSaleActive;
    /// @notice Returns true when the public sale is active, false otherwise.
    bool public publicSaleActive;

    /// @notice Keeps track of whether herolist has already minted or not. Max 1 mint.
    mapping(address => bool) public whitelistClaimed;
    /// @notice Keeps track of whether a public minter has already minted or not. Max 1 mint.
    mapping(address => bool) public publicClaimed;

    /// @notice Merkle root hash for whitelist verification.
    /// @dev Set to immutable instead of hard-coded to prevent human-error when deploying.
    bytes32 public merkleRoot;

    constructor(
        string memory _baseURI,
        bytes32 _merkleRoot
    ) ERC721("Chain Heroes", "HERO") Owned(msg.sender) {
        baseURI = _baseURI;
        merkleRoot = _merkleRoot;
        _balanceOf[msg.sender] = 10;
        unchecked {
            for (uint256 i = 0; i < 10; ++i) {
                _ownerOf[i] = msg.sender;
                emit Transfer(address(0), msg.sender, i);
            }
        }
    }

    /// @notice Allows the owner to change the base URI of CH's corresponding metadata.
    /// @param _uri The new URI to set the base URI to.
    function setURI(string calldata _uri) external onlyOwner {
        baseURI = _uri;
    }

    /// @notice The URI pointing to the metadata of a specific assett.
    /// @param _id The token ID of the requested hero. Hardcoded .json as suffix.
    /// @return The metadata URI.
    function tokenURI(
        uint256 _id
    ) public view override returns (string memory) {
        return string(abi.encodePacked(baseURI, _id.toString(), ".json"));
    }

    /// @notice Public Hero mint.
    /// @dev Allows any non-contract signer to mint a single Hero. Capped by 1.
    /// @dev Herolisted addresses can mint one during hero-sale and one during the public sale.
    /// @dev Current supply addition can be unchecked, as it cannot overflow.
    function publicMint() public payable {
        if (!publicSaleActive) revert PublicSaleNotActive();
        if (publicClaimed[msg.sender]) revert AlreadyClaimed();
        if (totalSupply + 1 > MAX_SUPPLY) revert MintedOut();
        if ((msg.value) < mintPrice) revert NotEnoughFunds();

        unchecked {
            publicClaimed[msg.sender] = true;
            _mint(msg.sender, totalSupply);
            ++totalSupply;
        }
    }

    /// @notice Mints a Hero for a signer on the herolist. Gets the tokenID correspondign to the current supply.
    /// @dev We do not keep track of the whitelist supply, considering only a total of 4337 addresses will be valid in the merkle tree.
    /// @dev This means that the maximum supply including full herolist mint and team mint can be 4347 at most, as each address can mint once.
    /// @dev Current supply addition can be unchecked, as it cannot overflow.
    /// @param _merkleProof The merkle proof based on the address of the signer as input.
    function heroListMint(bytes32[] calldata _merkleProof) public payable {
        if (!heroSaleActive) revert HeroSaleNotActive();
        if (whitelistClaimed[msg.sender]) revert AlreadyClaimed();
        if ((msg.value) < mintPrice) revert NotEnoughFunds();

        bytes32 leaf = keccak256(abi.encodePacked(msg.sender));
        if (!MerkleProof.verify(_merkleProof, merkleRoot, leaf))
            revert NotHeroListed();

        unchecked {
            whitelistClaimed[msg.sender] = true;
            _mint(msg.sender, totalSupply);
            ++totalSupply;
        }
    }

    /// @notice Flip the hero sale state.
    function flipHeroSaleState() public onlyOwner {
        heroSaleActive = !heroSaleActive;
    }

    /// @notice Flip the public sale state.
    function flipPublicSaleState() public onlyOwner {
        heroSaleActive = false;
        publicSaleActive = !publicSaleActive;
    }

    /// @notice Set the price of mint, in case there is no mint out.
    function setPrice(uint256 _targetPrice) public onlyOwner {
        mintPrice = _targetPrice;
    }

    /// @notice Transfer all funds from contract to the contract deployer address.
    function withdraw() public onlyOwner {
        (bool success, ) = msg.sender.call{value: address(this).balance}("");
        require(success);
    }

    /// @notice Set the merkle root.
    function setMerkleRoot(bytes32 _merkleRoot) public onlyOwner {
        merkleRoot = _merkleRoot;
    }
}
