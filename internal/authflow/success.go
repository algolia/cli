package authflow

const oauthSuccessPage = `
<!doctype html>
<meta charset="utf-8">
<title>Success: Algolia CLI</title>
<style type="text/css">
body {
  color: #1B1F23;
  background: #F6F8FA;
  font-size: 14px;
  font-family: -apple-system, "Segoe UI", Helvetica, Arial, sans-serif;
  line-height: 1.5;
  max-width: 620px;
  margin: 28px auto;
  text-align: center;
}

h1 {
  font-size: 24px;
  margin-bottom: 0;
}

p {
  margin-top: 0;
}

.box {
  border: 1px solid #E1E4E8;
  background: white;
  padding: 24px;
  margin: 28px;
}
</style>
<body>
  <div class="box">
    <h1>Successfully authenticated Algolia CLI</h1>
    <p>You may now close this tab and return to the terminal.</p>
  </div>
</body>
`
