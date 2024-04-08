// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.12;

error ZeroAddress();
error InvalidStartIndex();
error InvalidConfirmer();
error NotWhitelister();
error InvalidSender();
error InvalidReferenceBlockNum();
error InsufficientThreshold();
error InsufficientThresholdPercentages();
error InvalidQuorumParam();
error InvalidQuorumThresholdPercentage();
error AlreadyInAllowlist();
error NotInAllowlist();
error AlreadyAdded();
error ResolvedAlert();
error AlreadyEnabled();
error AlreadyDisabled();

// Common
error AlreadyInitialized();
error NotInitialized();
error ZeroValue();

error UselessAlert();
error InvalidAlert();
error InvalidAlertType();
error InvalidProvedIndex();
error InvalidCheckpoint();
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
