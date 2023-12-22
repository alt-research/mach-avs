// SPDX-License-Identifier: UNLICENSED
// SEE LICENSE IN https://files.altlayer.io/Alt-Research-License-1.md
// Copyright Alt Research Ltd. 2023. All rights reserved.
//
// You acknowledge and agree that Alt Research Ltd. ("Alt Research") (or Alt
// Research's licensors) own all legal rights, titles and interests in and to the
// work, software, application, source code, documentation and any other documents

pragma solidity =0.8.12;

// Common
error AlreadyInitialized();
error NotInitialized();
error ZeroAddress();
error ZeroValue();

error UselessAlert();
error InvalidAlert();
error InvalidAlertType();
error InvalidProvedIndex();
error InvalidPerCheckpoint();
error InvalidIndex();

error ProveImageIdMismatch();
error ProveBlockNumberMismatch();
error ProveOutputRootMismatch();
error ParentCheckpointNumberMismatch();
error ParentCheckpointOutputRootMismatch();
error ProveVerifyFailed();
error InvalidJournal();
error NoAlert();
error NotOperator();
