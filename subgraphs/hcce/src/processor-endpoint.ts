import { BigInt } from "@graphprotocol/graph-ts";
import {
  RequestSubmitted as RequestSubmittedEvent,
  RequestCompleted as RequestCompletedEvent,
  UserEvent as UserEventEvent,
  DeployRequestSubmitted as DeployRequestSubmittedEvent,
  DeployRequestCompleted as DeployRequestCompletedEvent,
} from "../generated/ProcessorEndpoint/ProcessorEndpoint";
import { RequestSubmitted, RequestCompleted, UserEvent, DeployRequestCompleted, DeployRequestSubmitted } from "../generated/schema";

const SORT_BASE = BigInt.fromI64(1000000000000);

export function handleRequestSubmitted(event: RequestSubmittedEvent): void {
  const id = event.transaction.hash.concatI32(event.logIndex.toI32()).toHex();
  let entity = new RequestSubmitted(id);

  entity.applicationId = event.params.applicationId;
  entity.requestId = event.params.requestId;
  entity.sender = event.params.sender;
  entity.blockNumber = event.block.number;
  entity.blockTimestamp = event.block.timestamp;

  entity.save();
}

export function handleRequestCompleted(event: RequestCompletedEvent): void {
  const id = event.transaction.hash.concatI32(event.logIndex.toI32()).toHex();
  let entity = new RequestCompleted(id);

  entity.applicationId = event.params.applicationId;
  entity.requestId = event.params.requestId;
  entity.applicationFees = event.params.applicationFees;
  entity.status = event.params.status;
  entity.errorCode = event.params.errorCode;
  entity.errorMessage = event.params.errorMessage;
  entity.blockNumber = event.block.number;
  entity.blockTimestamp = event.block.timestamp;

  entity.save();
}

export function handleUserEvent(event: UserEventEvent): void {
  const id = event.transaction.hash.concatI32(event.logIndex.toI32()).toHex();
  let entity = new UserEvent(id);

  entity.applicationId = event.params.applicationId;
  entity.requestId = event.params.requestId;
  entity.eventSubType = event.params.eventSubType;
  entity.encryptedData = event.params.encryptedData;
  entity.blockNumber = event.block.number;
  entity.logIndex = event.logIndex;
  entity.sortKey = event.block.number.times(SORT_BASE).plus(event.logIndex);
  entity.blockTimestamp = event.block.timestamp;

  entity.save();
}


export function handleDeployRequestSubmitted(event: DeployRequestSubmittedEvent): void {
  const id = event.transaction.hash.concatI32(event.logIndex.toI32()).toHex();
  let entity = new DeployRequestSubmitted(id);

  entity.applicationId = event.params.applicationId;
  entity.requestId = event.params.requestId;
  entity.sender = event.params.sender;
  entity.blockNumber = event.block.number;
  entity.blockTimestamp = event.block.timestamp;

  entity.save();
}

export function handleDeployRequestCompleted(event: DeployRequestCompletedEvent): void {
  const id = event.transaction.hash.concatI32(event.logIndex.toI32()).toHex();
  let entity = new DeployRequestCompleted(id);

  entity.applicationId = event.params.applicationId;
  entity.requestId = event.params.requestId;
  entity.applicationFees = event.params.applicationFees;
  entity.status = event.params.status;
  entity.errorCode = event.params.errorCode;
  entity.errorMessage = event.params.errorMessage;
  entity.blockNumber = event.block.number;
  entity.blockTimestamp = event.block.timestamp;

  entity.save();
}