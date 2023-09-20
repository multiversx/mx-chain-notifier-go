import json
import pytest

from receive import handle_txs_block_data


def test_handle_txs_block_data():
    with open('txs_test.json') as f:
        txs_data = json.load(f)

        alert = handle_txs_block_data(txs_data)
        assert alert == "ALERT: found tx 6b100923f6a3e29c14d8b75910bfa2e44876e77de15c1e2a9d65cfacdab8d59d for erd1qqqqqqqqqqqqqpgqvrdx4n97ensz534rw6s2pn5rfugwv7sl0tgqwzfgxv with ChangeOwnerAddress"
