// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import '@openzeppelin/contracts/access/AccessControl.sol';
import '@openzeppelin/contracts/utils/ReentrancyGuard.sol';
import '@openzeppelin/contracts/token/ERC20/IERC20.sol';
import '@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol';
import '@openzeppelin/contracts/token/ERC20/extensions/IERC20Permit.sol';
import '@openzeppelin/contracts/utils/cryptography/ECDSA.sol';
import '@openzeppelin/contracts/utils/cryptography/EIP712.sol';
import '@openzeppelin/contracts/utils/Strings.sol';

import './interfaces/ITeeAuthenticator.sol';
import './interfaces/IProcessorEndpoint.sol';
import './interfaces/IAuthorityRegistry.sol';
import './TokenAllowlist.sol';
import './Structs.sol';

/// @title ProcessorEndpoint
/// @notice Implementation of the processor endpoint interface.
contract ProcessorEndpoint is TokenAllowlist, IProcessorEndpoint, ReentrancyGuard, EIP712 {
  using SafeERC20 for IERC20;

  //constants
  bytes32 public constant UPDATE_STATUS_ROLE = keccak256('UPDATE_STATUS_ROLE');
  bytes32 public constant ADMIN = keccak256('ADMIN');
  bytes32 public constant DEPLOYER_ROLE = keccak256('DEPLOYER_ROLE');
  uint8 public constant PROTOCOL_VERSION = 0;
  //state variables
  mapping(uint64 => bytes32) public applicationStateRoots;
  uint256 public maxNumOfApplications = 10;
  uint256 public availableDeploySlots = maxNumOfApplications;

  mapping(bytes32 => Structs.PendingRequest) public requestById;
  mapping(uint256 => bytes32) private _requestIdByOrder;
  uint256 private _head;
  uint256 private _tail;
  uint256 public maxQueueSize = 10;

  ITeeAuthenticator public teeAuthenticator;
  IAuthorityRegistry public authorityRegistry;

  // Pull payment pattern state — per-token, per-payee
  mapping(address => mapping(address => uint256)) public pendingClaims;
  mapping(address => uint256) public totalPendingClaims;

  // Per-app, per-token custody tracking for solvency isolation.
  // Credited on submitRequest (assetAmount only; fees are tracked globally),
  // debited on stateUpdate: success path (withdrawals), error path (assetAmount).
  mapping(uint64 => mapping(address => uint256)) public appCustody;
  mapping(address => uint256) public totalAppCustody;

  uint256 public minFeePerRequest;
  address payable public feeCollector;

  // EIP-712 typehash for facilitator request authorization
  bytes32 public constant REQUEST_AUTHORIZATION_TYPEHASH =
    keccak256(
      'RequestAuthorization(address sender,uint8 protocolVersion,uint64 applicationId,uint8 requestType,bytes32 payloadHash,address tokenAddress,uint256 assetAmount,uint256 nonce,uint256 deadline)'
    );

  // Sequential nonces per user for facilitator replay protection
  mapping(address => uint256) public facilitatorNonces;

  modifier validProtocolVersion(uint8 protocolVersion) {
    if (protocolVersion != PROTOCOL_VERSION) revert InvalidProtocolVersion();
    _;
  }

  modifier validApplicationId(uint64 applicationId) {
    if (applicationStateRoots[applicationId] == bytes32(0)) revert InvalidApplicationId();
    _;
  }

  /// @param _teeAuthenticator Contract used to verify update signatures.
  /// @param _authorityRegistry Registry for authority checks.
  /// @param updateStatusOperator Initial operator for status updates.
  /// @param admin Initial admin address.
  /// @param _minFeePerRequest Minimum fee enforced per request.
  constructor(
    ITeeAuthenticator _teeAuthenticator,
    IAuthorityRegistry _authorityRegistry,
    address updateStatusOperator,
    address admin,
    uint256 _minFeePerRequest
  ) EIP712('Vela', Strings.toString(PROTOCOL_VERSION)) {
    if (
      address(_teeAuthenticator) == address(0) ||
      address(_authorityRegistry) == address(0) ||
      updateStatusOperator == address(0) ||
      admin == address(0)
    ) revert AddressCantBeZero();

    teeAuthenticator = _teeAuthenticator;
    authorityRegistry = _authorityRegistry;
    feeCollector = payable(updateStatusOperator);
    _grantRole(UPDATE_STATUS_ROLE, updateStatusOperator);
    _grantRole(ADMIN, admin);
    _setRoleAdmin(DEPLOYER_ROLE, ADMIN);
    _grantRole(DEPLOYER_ROLE, admin);
    minFeePerRequest = _minFeePerRequest;
  }

  /// @inheritdoc IProcessorEndpoint
  function submitRequest(
    uint8 protocolVersion,
    uint64 applicationId,
    Structs.RequestType requestType,
    bytes calldata payload,
    address tokenAddress,
    uint256 assetAmount,
    uint256 maxFeeValue // part of the sent value reserved for fee payment
  )
    external
    payable
    validProtocolVersion(protocolVersion)
    validApplicationId(applicationId)
    nonReentrant
    returns (bytes32)
  {
    if (requestType == Structs.RequestType.DEPLOYAPP) revert InvalidRequestType();

    //check values
    if (maxFeeValue < minFeePerRequest) revert FeeValueBelowMinimum();

    if (tokenAddress == ETH_TOKEN) {
      if (msg.value != assetAmount + maxFeeValue) revert InvalidValue();
    } else {
      if (msg.value != maxFeeValue) revert InvalidValue();
      if (assetAmount == 0) revert InvalidValue();
      if (!allowedTokens[tokenAddress]) revert TokenNotAllowed();
      // Pull ERC-20 tokens with balance-before/after check
      IERC20 token = IERC20(tokenAddress);
      uint256 balanceBefore = token.balanceOf(address(this));
      token.safeTransferFrom(msg.sender, address(this), assetAmount);
      uint256 received = token.balanceOf(address(this)) - balanceBefore;
      if (received != assetAmount) revert TransferAmountMismatch();
    }

    //check queue size
    if (getPendingRequestsSize() >= maxQueueSize) revert QueueThresholdExceeded();

    if (requestType == Structs.RequestType.ASSOCIATEKEY) {
      //if requestype is associatekey, the payload must be 133 bytes (key only) or 226 bytes (key + encrypted seed)
      if (payload.length != 133 && payload.length != 226) revert InvalidPayload();
    } else if (requestType == Structs.RequestType.DEANONYMIZATION) {
      // only allowed authorities can request deanonymization
      if (!authorityRegistry.checkAuthorityIsAllowed(applicationId, msg.sender)) {
        revert AuthorityNotAllowed();
      }
    }

    //create request
    bytes32 requestId = generateRequestId(
      msg.sender,
      applicationId,
      requestType,
      payload,
      tokenAddress,
      assetAmount,
      _tail
    );
    requestById[requestId] = Structs.PendingRequest({
      timestamp: block.timestamp,
      tokenAddress: tokenAddress,
      assetAmount: assetAmount,
      maxFeeValue: maxFeeValue,
      requestId: requestId,
      payload: payload,
      sender: msg.sender,
      facilitator: address(0),
      applicationId: applicationId,
      protocolVersion: protocolVersion,
      requestType: requestType
    });
    _requestIdByOrder[_tail] = requestId;
    appCustody[applicationId][tokenAddress] += assetAmount;
    totalAppCustody[tokenAddress] += assetAmount;

    unchecked {
      ++_tail;
    }

    //emit event
    emit RequestSubmitted(applicationId, requestId, msg.sender, address(0));

    return requestId;
  }

  /// @inheritdoc IProcessorEndpoint
  function submitRequestFor(
    address sender,
    uint8 protocolVersion,
    uint64 applicationId,
    Structs.RequestType requestType,
    bytes calldata payload,
    address tokenAddress,
    uint256 assetAmount,
    uint256 deadline,
    bytes calldata requestSignature,
    bytes calldata depositPermit
  )
    external
    payable
    validProtocolVersion(protocolVersion)
    validApplicationId(applicationId)
    nonReentrant
    returns (bytes32)
  {
    // 1. Only ASSOCIATEKEY and PROCESS are supported
    if (
      requestType != Structs.RequestType.ASSOCIATEKEY &&
      requestType != Structs.RequestType.PROCESS
    ) revert InvalidRequestType();

    // 2. Verify deadline not expired
    if (block.timestamp > deadline) revert DeadlineExpired();

    // 3. Read current nonce and build EIP-712 hash
    uint256 nonce = facilitatorNonces[sender];
    bytes32 structHash = keccak256(
      abi.encode(
        REQUEST_AUTHORIZATION_TYPEHASH,
        sender,
        protocolVersion,
        applicationId,
        uint8(requestType),
        keccak256(payload),
        tokenAddress,
        assetAmount,
        nonce,
        deadline
      )
    );
    bytes32 digest = _hashTypedDataV4(structHash);

    // 4. Recover user address from EIP-712 request signature and verify
    address recoveredSigner = ECDSA.recover(digest, requestSignature);
    if (recoveredSigner == address(0)) revert InvalidSignature();
    if (recoveredSigner != sender) revert InvalidSigner();

    // 5. Consume nonce (replay protection)
    facilitatorNonces[sender] = nonce + 1;

    // 6. Validate fee
    uint256 maxFeeValue = msg.value;
    if (maxFeeValue < minFeePerRequest) revert FeeValueBelowMinimum();

    // 7. Validate token and handle deposit
    if (assetAmount > 0) {
      if (tokenAddress == ETH_TOKEN) revert InvalidValue();
      if (!allowedTokens[tokenAddress]) revert TokenNotAllowed();

      // Decode deposit permit and execute EIP-2612 permit + transferFrom
      if (depositPermit.length != 96) revert InvalidPermit();
      (uint8 v, bytes32 r, bytes32 s) = abi.decode(depositPermit, (uint8, bytes32, bytes32));

      // Check current allowance before calling permit
      IERC20 token = IERC20(tokenAddress);
      if (token.allowance(sender, address(this)) < assetAmount) {
        IERC20Permit(tokenAddress).permit(sender, address(this), assetAmount, deadline, v, r, s);
      }

      uint256 balanceBefore = token.balanceOf(address(this));
      token.safeTransferFrom(sender, address(this), assetAmount);
      uint256 received = token.balanceOf(address(this)) - balanceBefore;
      if (received != assetAmount) revert TransferAmountMismatch();

      appCustody[applicationId][tokenAddress] += assetAmount;
      totalAppCustody[tokenAddress] += assetAmount;
    }

    // 8. Check queue size
    if (getPendingRequestsSize() >= maxQueueSize) revert QueueThresholdExceeded();

    // 9. Validate payload for ASSOCIATEKEY
    if (requestType == Structs.RequestType.ASSOCIATEKEY) {
      if (payload.length != 133 && payload.length != 226) revert InvalidPayload();
    }

    // 10. Create PendingRequest with sender = user (not msg.sender)
    bytes32 requestId = generateRequestId(
      sender,
      applicationId,
      requestType,
      payload,
      tokenAddress,
      assetAmount,
      _tail
    );
    requestById[requestId] = Structs.PendingRequest({
      timestamp: block.timestamp,
      tokenAddress: tokenAddress,
      assetAmount: assetAmount,
      maxFeeValue: maxFeeValue,
      requestId: requestId,
      payload: payload,
      sender: sender,
      facilitator: msg.sender,
      applicationId: applicationId,
      protocolVersion: protocolVersion,
      requestType: requestType
    });
    _requestIdByOrder[_tail] = requestId;

    unchecked {
      ++_tail;
    }

    // 11. Emit event with user address
    emit RequestSubmitted(applicationId, requestId, sender, msg.sender);

    return requestId;
  }

  /// @inheritdoc IProcessorEndpoint
  function getFacilitatorNonce(address user) external view returns (uint256) {
    return facilitatorNonces[user];
  }

  /// @inheritdoc IProcessorEndpoint
  function submitDeployRequest(
    uint8 protocolVersion,
    bytes calldata payload
  ) external payable validProtocolVersion(protocolVersion) nonReentrant returns (bytes32) {
    if (!hasRole(DEPLOYER_ROLE, msg.sender)) revert DeployerNotAllowed();
    if (availableDeploySlots == 0) revert MaxNumOfApplicationsExceeded();
    //check queue size
    if (getPendingRequestsSize() >= maxQueueSize) revert QueueThresholdExceeded();
    if (msg.value < minFeePerRequest) revert FeeValueBelowMinimum();

    --availableDeploySlots;

    Structs.RequestType requestType = Structs.RequestType.DEPLOYAPP;
    //create request
    bytes32 requestId = generateRequestId(
      msg.sender,
      0, // deploy requests have applicationId 0, a unique applicationId will be derived from the requestId for each deploy request to avoid collisions with regular requests and to group deploy requests together
      requestType,
      payload,
      ETH_TOKEN,
      0,
      _tail
    );

    uint64 applicationId = uint64(bytes8(requestId)); // Derive a unique application ID from the request ID for deploy requests
    requestById[requestId] = Structs.PendingRequest({
      timestamp: block.timestamp,
      tokenAddress: ETH_TOKEN,
      assetAmount: 0,
      maxFeeValue: msg.value,
      requestId: requestId,
      payload: payload,
      sender: msg.sender,
      facilitator: address(0),
      applicationId: applicationId,
      protocolVersion: protocolVersion,
      requestType: requestType
    });
    _requestIdByOrder[_tail] = requestId;

    unchecked {
      ++_tail;
    }

    //emit event
    emit DeployRequestSubmitted(applicationId, requestId, msg.sender);

    return requestId;
  }

  function _removeRequest() private {
    delete requestById[_requestIdByOrder[_head]];
    delete _requestIdByOrder[_head];
    unchecked {
      ++_head;
    }
  }

  function _markRequestCompleted(
    uint64 applicationId,
    bytes32 requestId,
    uint256 applicationFees,
    Structs.RequestResult result,
    Structs.ErrorCode errCode,
    string memory errorMsg,
    Structs.RequestType requestType
  ) private {
    _removeRequest();

    if (requestType == Structs.RequestType.DEPLOYAPP) {
      emit DeployRequestCompleted(
        applicationId,
        requestId,
        applicationFees,
        result,
        errCode,
        errorMsg
      );
    } else {
      emit RequestCompleted(applicationId, requestId, applicationFees, result, errCode, errorMsg);
    }

    _asyncTransfer(ETH_TOKEN, feeCollector, applicationFees);
  }

  /// @inheritdoc IProcessorEndpoint
  function getPendingRequestsSize() public view returns (uint256) {
    if (_tail > _head) {
      return (_tail - _head);
    } else {
      return 0;
    }
  }

  /// @inheritdoc IProcessorEndpoint
  function getPendingRequests() external view returns (Structs.PendingRequest[] memory) {
    uint256 numOfPendingRequests = getPendingRequestsSize();

    Structs.PendingRequest[] memory res = new Structs.PendingRequest[](numOfPendingRequests);
    uint256 i = _head;
    uint256 tail = _tail;
    uint256 j;
    while (i < tail) {
      bytes32 requestId = _requestIdByOrder[i];
      res[j] = requestById[requestId];
      unchecked {
        ++i;
        ++j;
      }
    }

    return res;
  }

  /// @inheritdoc IProcessorEndpoint
  function getPendingRequestsPage(
    uint256 offset,
    uint256 limit
  ) external view returns (Structs.PendingRequest[] memory) {
    uint256 size = getPendingRequestsSize();
    if (offset >= size || limit == 0) {
      return new Structs.PendingRequest[](0);
    }

    uint256 end = offset + limit;
    if (end > size) {
      end = size;
    }
    uint256 count = end - offset;

    Structs.PendingRequest[] memory res = new Structs.PendingRequest[](count);
    uint256 i = _head + offset;
    uint256 stop = _head + end;
    uint256 j;
    while (i < stop) {
      bytes32 requestId = _requestIdByOrder[i];
      res[j] = requestById[requestId];
      unchecked {
        ++i;
        ++j;
      }
    }

    return res;
  }

  //update status
  /// @inheritdoc IProcessorEndpoint
  function stateUpdate(
    uint64 applicationId,
    bytes32 prevStateRoot,
    bytes32 newStateRoot,
    bytes32 processedRequestId,
    bytes[] calldata events,
    string[] calldata eventSubTypes,
    Structs.WithdrawalRequest[] calldata withdrawalRequests,
    uint256 refund,
    uint256 applicationFees,
    Structs.ErrorCode errorCode,
    string calldata errorMsg,
    bytes calldata signature
  ) external onlyRole(UPDATE_STATUS_ROLE) nonReentrant {
    //check valid request
    if (!isCurrentPendingRequest(processedRequestId)) revert InvalidRequestId();

    // Check application Id
    Structs.PendingRequest storage requestInfo = requestById[processedRequestId];
    if (applicationId != requestInfo.applicationId) revert InvalidApplicationId();

    //check prev state root
    if (prevStateRoot != applicationStateRoots[applicationId]) revert InvalidStateRoot();

    uint256 eventsLength = events.length;
    uint256 eventSubTypesLength = eventSubTypes.length;
    if (eventsLength != eventSubTypesLength) revert InvalidPayload();

    //check signature
    Structs.SignatureParams memory sigParams = Structs.SignatureParams({
      applicationId: applicationId,
      prevStateRoot: prevStateRoot,
      newStateRoot: newStateRoot,
      processedRequestId: processedRequestId,
      events: events,
      eventSubTypes: eventSubTypes,
      withdrawalRequests: withdrawalRequests,
      refundAmount: refund,
      applicationFee: applicationFees,
      errorCode: errorCode,
      errorMsg: errorMsg
    });
    if (!teeAuthenticator.checkSignature(sigParams, signature)) revert InvalidSignature();

    //check values

    uint256 maxFeeValue = requestInfo.maxFeeValue;
    address payable sender = payable(requestInfo.sender);
    // Fee refund recipient: facilitator if present, otherwise sender
    address payable feeRecipient = requestInfo.facilitator != address(0)
      ? payable(requestInfo.facilitator)
      : sender;

    // Handle error case (signed error payload from TEE)
    if (errorCode != Structs.ErrorCode.NO_ERROR) {
      // For errors: state unchanged (prevStateRoot == newStateRoot), no events, no withdrawals
      // Refund user (minus minimum fee) and collect minimum fee
      if (eventsLength != 0 || withdrawalRequests.length != 0) revert InvalidPayload();
      if (applicationStateRoots[applicationId] != newStateRoot) revert InvalidStateRoot();

      // Per-app per-token solvency check, then ETH balance check for fee outflow.
      uint256 assetAmount = requestInfo.assetAmount;
      address reqTokenAddress = requestInfo.tokenAddress;
      if (assetAmount > appCustody[applicationId][reqTokenAddress]) revert InsufficientAppBalance();
      if (requestInfo.maxFeeValue > _getAvailableEthBalance()) revert InsufficientBalance();
      appCustody[applicationId][reqTokenAddress] -= assetAmount;
      totalAppCustody[reqTokenAddress] -= assetAmount;

      // Refund business-asset deposit in its original token to the user.
      // Fee refund is always in ETH, routed to facilitator if present.
      uint256 feeRefund = requestInfo.maxFeeValue - minFeePerRequest;
      if (reqTokenAddress == ETH_TOKEN) {
        // For ETH requests (only direct path, never facilitated)
        uint256 totalRefund = assetAmount + feeRefund;
        if (totalRefund > 0) {
          _asyncTransfer(ETH_TOKEN, sender, totalRefund);
          emit Refund(applicationId, processedRequestId, ETH_TOKEN, sender, totalRefund);
        }
      } else {
        // For ERC-20 requests, asset refund in token to user, fee refund in ETH to feeRecipient
        if (assetAmount > 0) {
          _asyncTransfer(reqTokenAddress, sender, assetAmount);
          emit Refund(applicationId, processedRequestId, reqTokenAddress, sender, assetAmount);
        }
        if (feeRefund > 0) {
          _asyncTransfer(ETH_TOKEN, feeRecipient, feeRefund);
          emit Refund(applicationId, processedRequestId, ETH_TOKEN, feeRecipient, feeRefund);
        }
      }

      if (requestInfo.requestType == Structs.RequestType.DEPLOYAPP) {
        unchecked {
          ++availableDeploySlots;
        }
      }

      _markRequestCompleted(
        applicationId,
        processedRequestId,
        minFeePerRequest,
        Structs.RequestResult.FAILED,
        Structs.ErrorCode(errorCode),
        errorMsg,
        requestInfo.requestType
      );

      return;
    }

    // Handle success case
    // State cannot remain the same
    if (applicationStateRoots[applicationId] == newStateRoot) revert InvalidStateRoot();

    if (refund + applicationFees != maxFeeValue) revert InvalidValue();
    if (applicationFees < minFeePerRequest) {
      revert InvalidValue();
    }

    //check withdrawal sums and debit per-app per-token custody
    uint256 i;
    uint256 withdrawalsLength = withdrawalRequests.length;
    uint256 ethWithdrawalSum;
    while (i < withdrawalsLength) {
      address wToken = withdrawalRequests[i].tokenAddress;
      uint256 wAmount = withdrawalRequests[i].amount;
      if (wAmount > appCustody[applicationId][wToken]) revert InsufficientAppBalance();
      appCustody[applicationId][wToken] -= wAmount;
      totalAppCustody[wToken] -= wAmount;
      if (wToken == ETH_TOKEN) {
        ethWithdrawalSum += wAmount;
      } else {
        // Per-token solvency: contract must hold enough ERC-20 to cover all obligations
        if (
          IERC20(wToken).balanceOf(address(this)) <
          totalAppCustody[wToken] + totalPendingClaims[wToken] + wAmount
        ) revert InsufficientBalance();
      }
      unchecked {
        ++i;
      }
    }

    // ETH solvency: contract must hold enough ETH to cover fee outflow + ETH withdrawals
    uint256 totalEthOutflow = refund + applicationFees + ethWithdrawalSum;
    if (totalEthOutflow > _getAvailableEthBalance()) revert InsufficientBalance();

    //emit encrypted event
    i = 0;
    while (i < eventsLength) {
      emit UserEvent(applicationId, processedRequestId, eventSubTypes[i], events[i]);
      unchecked {
        ++i;
      }
    }

    Structs.RequestType reqType = requestInfo.requestType;
    if (reqType == Structs.RequestType.DEANONYMIZATION) {
      //a completed DEANONYMIZATION request must have always generated a report
      emit ReportGenerated(applicationId, processedRequestId);
    }

    //update state root and request
    applicationStateRoots[applicationId] = newStateRoot;
    emit StateRootUpdate(applicationId, processedRequestId, prevStateRoot, newStateRoot);

    //credit refund to feeRecipient's pending balance (pull pattern) — refund is always ETH
    if (refund > 0) {
      _asyncTransfer(ETH_TOKEN, feeRecipient, refund);
      emit Refund(applicationId, processedRequestId, ETH_TOKEN, feeRecipient, refund);
    }

    //credit withdrawals to receivers' pending balances
    i = 0;
    while (i < withdrawalsLength) {
      _asyncTransfer(
        withdrawalRequests[i].tokenAddress,
        withdrawalRequests[i].receiver,
        withdrawalRequests[i].amount
      );
      emit Withdrawal(
        applicationId,
        processedRequestId,
        withdrawalRequests[i].tokenAddress,
        withdrawalRequests[i].receiver,
        withdrawalRequests[i].amount
      );
      unchecked {
        ++i;
      }
    }

    //set requests as completed
    _markRequestCompleted(
      applicationId,
      processedRequestId,
      applicationFees,
      Structs.RequestResult.COMPLETED,
      Structs.ErrorCode.NO_ERROR,
      '',
      reqType
    );
  }

  /// @inheritdoc IProcessorEndpoint
  function updateQueueThreshold(uint256 newThreshold) external onlyRole(ADMIN) {
    if (newThreshold == 0) revert InvalidValue();
    maxQueueSize = newThreshold;
    emit QueueThresholdUpdated(newThreshold);
  }

  /// @inheritdoc IProcessorEndpoint
  function updateMaxNumOfApplications(uint256 newMax) external onlyRole(ADMIN) {
    if (newMax == 0) revert InvalidValue();
    uint256 deployedApps = maxNumOfApplications - availableDeploySlots;
    if (newMax < deployedApps) revert InvalidValue();
    uint256 oldMax = maxNumOfApplications;
    maxNumOfApplications = newMax;
    availableDeploySlots = newMax - deployedApps;
    emit MaxNumberOfAppUpdated(oldMax, newMax);
  }

  /// @inheritdoc IProcessorEndpoint
  function updateFeeCollector(address payable newFeeCollector) external onlyRole(ADMIN) {
    if (newFeeCollector == address(0)) revert AddressCantBeZero();
    feeCollector = newFeeCollector;
    emit FeeCollectorUpdated(newFeeCollector);
  }

  /// @inheritdoc IProcessorEndpoint
  function addAllowedDeployer(address deployer) external onlyRole(ADMIN) {
    if (deployer == address(0)) revert AddressCantBeZero();
    _grantRole(DEPLOYER_ROLE, deployer);
  }

  /// @inheritdoc IProcessorEndpoint
  function removeAllowedDeployer(address deployer) external onlyRole(ADMIN) {
    if (deployer == address(0)) revert AddressCantBeZero();
    _revokeRole(DEPLOYER_ROLE, deployer);
  }

  /// @inheritdoc IProcessorEndpoint
  function isAllowedDeployer(address deployer) external view returns (bool allowed) {
    return hasRole(DEPLOYER_ROLE, deployer);
  }

  /// @inheritdoc IProcessorEndpoint
  function getNextPendingRequest()
    external
    view
    returns (Structs.PendingRequest memory, bytes32, bool success)
  {
    uint256 numOfRequests = getPendingRequestsSize();
    if (numOfRequests > 0) {
      bytes32 requestId = _requestIdByOrder[_head];
      Structs.PendingRequest storage req = requestById[requestId];
      return (req, applicationStateRoots[req.applicationId], true);
    }

    Structs.PendingRequest memory emptyReq;
    return (emptyReq, bytes32(0), false);
  }

  /// @inheritdoc IProcessorEndpoint
  function isCurrentPendingRequest(bytes32 requestId) public view returns (bool) {
    return getPendingRequestsSize() > 0 && _requestIdByOrder[_head] == requestId;
  }

  // Pull payment pattern functions
  function _asyncTransfer(address tokenAddress, address dest, uint256 amount) internal {
    pendingClaims[tokenAddress][dest] += amount;
    totalPendingClaims[tokenAddress] += amount;
  }

  function _getAvailableEthBalance() internal view returns (uint256) {
    return address(this).balance - totalPendingClaims[ETH_TOKEN] - totalAppCustody[ETH_TOKEN];
  }

  /// @inheritdoc IProcessorEndpoint
  function claim(address tokenAddress, address payable payee) public nonReentrant {
    uint256 amount = pendingClaims[tokenAddress][payee];
    if (amount == 0) return;

    pendingClaims[tokenAddress][payee] = 0;
    totalPendingClaims[tokenAddress] -= amount;

    emit PaymentWithdrawn(tokenAddress, payee, amount);

    if (tokenAddress == ETH_TOKEN) {
      (bool success, ) = payee.call{value: amount}('');
      if (!success) revert TransferFailed();
    } else {
      IERC20(tokenAddress).safeTransfer(payee, amount);
    }
  }

  /// @inheritdoc IProcessorEndpoint
  function generateRequestId(
    address sender,
    uint64 applicationId,
    Structs.RequestType requestType,
    bytes calldata payload,
    address tokenAddress,
    uint256 assetAmount,
    uint256 idx
  ) public pure returns (bytes32) {
    return keccak256(
      abi.encode(sender, applicationId, requestType, payload, tokenAddress, assetAmount, idx)
    );
  }
}
