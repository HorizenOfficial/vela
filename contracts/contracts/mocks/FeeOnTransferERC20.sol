// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import '@openzeppelin/contracts/token/ERC20/ERC20.sol';

/// @dev ERC20 that silently takes a 1-token fee on every transferFrom.
contract FeeOnTransferERC20 is ERC20 {
  constructor() ERC20('FeeToken', 'FEE') {}

  function mint(address to, uint256 amount) external {
    _mint(to, amount);
  }

  function transferFrom(address from, address to, uint256 amount) public override returns (bool) {
    _spendAllowance(from, _msgSender(), amount);
    // Transfer 1 less than requested to simulate fee-on-transfer
    _transfer(from, to, amount - 1);
    return true;
  }
}
