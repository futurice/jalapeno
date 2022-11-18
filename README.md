![Jalapeno - Empowered by Futurice](/docs/static/img/logo.png)

# Jalapeño

CLI for creating, managing and sharing spiced up project templates aka _recipes_

## Testing

### Debugging tests

To debug a single godog feature test in an IDE, in the `cmd` directory run

```shell
dlv test . --headless --listen 127.0.0.1:52800 -- --test.run '^TestFeatures$/^Foobar'
```

Where "Foobar" is the beginning of your feature description. Then connect to the running debugger from your IDE.
