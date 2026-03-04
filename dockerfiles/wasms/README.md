Legacy helper folder.

With WASM deploy v1, deploy no longer reads this folder directly.
Use wallet deploy command instead:

```
novaw deployapp --wasm /path/to/app.wasm --max-value-fee "100 wei"
```
