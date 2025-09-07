// SPDX-License-Identifier: MIT

pragma solidity ^0.8.26;

contract Store {
    event ItemSet(int256 key, int256 value);

    string public version;
    mapping (int256 => int256) public items;

    constructor(string memory _version) {
        version = _version;
    }

    function setItem(int256 key, int256 value) external {
        items[key] = value;
        emit ItemSet(key, value);
    }

}