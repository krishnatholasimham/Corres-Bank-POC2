# Enabling PBFT Consensus

## Using Consensus Plugin
Consensus requires some specific configuration and must be run on the latest version of the [hyperledger fabric](https://github.com/ANZ-Blockchain-Lab/fabric/commit/85768f2):

1. In `fabric/peer/core.yaml`, set the `peer.validator.consensus` value to `pbft`.
2. In `core.yaml`, make sure the `peer.id` is set sequentially as `vpX` where `X` is an integer that starts from `0` and goes to `N-1`. For example, with 4 validating peers, set the `peer.id` to`vp0`, `vp1`, `vp2`, `vp3`.
3. In `consensus/pbft/config.yaml` set the `general.mode` value to `batch` and the `general.N` value to the number of validating peers on the network, also set `general.batchsize` to the number of transactions per batch.
4. In `consensus/pbft/config.yaml`, optionally set timer values for the batch period (`general.timeout.batch`), the acceptable delay between request and execution (`general.timeout.request`), and for view-change (`general.timeout.viewchange`)

See `core.yaml` and `consensus/pbft/config.yaml` for more detail.

All of these setting may be overriden via the command line environment variables, eg. `CORE_PEER_VALIDATOR_CONSENSUS_PLUGIN=pbft` or `CORE_PBFT_GENERAL_MODE=batch`


## Naming the validators

When using the `pbft` module, make sure the validators are named sequentially using the `vpX` naming scheme, where `X` is an integer that starts from `0` and goes to `N-1`. For example, with 4 validating peers, you would set the `peer.id` keys on your validators to `vp0`, `vp1`, `vp2`, and `vp3`. [For those wondering why we do this: every validator in PBFT needs to maintain the same sorted list of validators, so that in a view change from `v` to `v+1`, all validators point towards the same new primary. Until whitelisting is implemented (see "Roadmap" section below), the `vpX` naming scheme is the most effective --although admittedly flaky and arbitrary-- way of doing this on the fly.]
