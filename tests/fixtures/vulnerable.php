<?php
// Critical
eval($_GET['code']);
exec($_POST['cmd']);
system($_REQUEST['cmd']);
passthru($_GET['p']);

// SQLi
$query = "SELECT * FROM users WHERE id = " . $_GET['id'];
$query2 = "DELETE FROM users WHERE username = '" . $_POST['user'] . "'";

// Redirects
header("Location: " . $_GET['url']);
header("Location: https://example.com/" . $_REQUEST['page']);

// Warning
$hash = md5('password123');
$form = '<form method="post" action="/submit">';

// XSS
?>
<script>
document.write("<img src=x onerror=alert(1)>");  // reflected XSS
location.href = userInput + "&next";             // JS redirect
innerHTML = "<div>" + userInput + "</div>";      // DOM XSS
</script>
<?php
// More RCE
shell_exec($_GET['cmd']);
?>