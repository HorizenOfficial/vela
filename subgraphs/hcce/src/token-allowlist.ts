import {
  TokenAllowed as TokenAllowedEvent,
  TokenRemoved as TokenRemovedEvent,
} from "../generated/TokenAllowlist/TokenAllowlist";
import { TokenAllowed, TokenRemoved } from "../generated/schema";

export function handleTokenAllowed(event: TokenAllowedEvent): void {
  let entity = new TokenAllowed(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  );

  entity.token = event.params.token;
  entity.blockNumber = event.block.number;
  entity.blockTimestamp = event.block.timestamp;

  entity.save();
}

export function handleTokenRemoved(event: TokenRemovedEvent): void {
  let entity = new TokenRemoved(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  );

  entity.token = event.params.token;
  entity.blockNumber = event.block.number;
  entity.blockTimestamp = event.block.timestamp;

  entity.save();
}
