// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {ClassInput} from "../ClassInput.sol";

struct MsgNewClass {
    address creator;
    ClassInput input;
}
