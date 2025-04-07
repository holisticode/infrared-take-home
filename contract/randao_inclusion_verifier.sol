pragma solidity ^0.8.29;

library MerklePath{
  function verify(
    bytes32 leaf,
    uint64 index,
    bytes32[] memory branch,
    bytes32 root
  ) internal pure returns (bool) {
    bytes32 value = leaf;
    for (uint256 i = 0; i< branch.length; i++) {
      if ((index >> i) & 1 == 0 ) {
        value = sha256(abi.encodePacked(value, branch[i]));
      } else {
        value = sha256(abi.encodePacked(branch[i], value));
      }
    }
    return value == root;
  }
}

contract RandaoVerifier {
    address beacon_root_contract = 0x000F3df6D732807Ef1319fB7B8bB8522d0Beac02;

    struct BeaconBlockHeader {
      uint64 slot;
      uint64 proposer_index;
      bytes32 parent_root;
      bytes32 state_root;
      bytes32 body_root;
    }

    function getBeaconRoot(uint256 timestamp) public returns (bytes32) {
        (bool ret, bytes memory data) = beacon_root_contract.call(bytes.concat(bytes32(timestamp)));
        require(ret);
        return bytes32(data);
    }

    function verifyRandaoMix(
      uint256 timestamp,
      bytes32 randaoMix,
      uint64 leafIndex,
      bytes32[] calldata proof,
      bytes32 blockHeaderRoot,
      bytes32 stateRoot 
    ) external returns (bool) {
      // check the provided block root is indeed the one for this timestamp
      if (getBeaconRoot(timestamp) != blockHeaderRoot) return false;

      bool proofValid = MerklePath.verify(randaoMix,leafIndex,proof,stateRoot);
      return proofValid;
      
    }
}
