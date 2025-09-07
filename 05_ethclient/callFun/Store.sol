// SPDX-License-Identifier: MIT

pragma solidity ^0.8.26;

contract Store {
    event ItemSet(address indexed  address1, int256 indexed key1, string indexed keyStr, int256 key, int256 value);
    string public version;
    mapping(int256 => int256) public items;

    constructor(string memory _version) {
        version = _version;
    }
    
    function setItem(int256 key, int256 value) external {
        items[key] = value;
        emit ItemSet(msg.sender, key, version, key, value);
    }

}