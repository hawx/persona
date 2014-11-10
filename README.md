# persona

Package for using [persona][] in webapps.

``` go
store := persona.NewStore(cookieSecret)
persona := persona.New(store, "http://localhost:8080", []string{"john.doe@example.com"})

http.Handle("/", persona.Switch(SignInHandler, SignedOutHandler))
http.Handle("/secret", persona.Protect(SecretHandler))
http.Handle("/sign-in", persona.SignIn)
http.Handle("/sign-out", persona.SignOut)
```

then add this to your html,

``` html
<script src="http://code.jquery.com/jquery-2.1.1.min.js"></script>
<script src="https://login.persona.org/include.js"></script>

<script>
  function gotAssertion(assertion) {
    // got an assertion, now send it up to the server for verification
    if (assertion !== null) {
      $.ajax({
        type: 'POST',
        url: '/sign-in',
        data: { assertion: assertion },
        success: function(res, status, xhr) {
          window.location.reload();
        },
        error: function(xhr, status, res) {
          alert("sign-in failure" + res);
        }
     });
    }
  };

  jQuery(function($) {
    $('#browserid').click(function() {
      navigator.id.get(gotAssertion);
    });
  });
</script>
```

[persona]: https://login.persona.org/
