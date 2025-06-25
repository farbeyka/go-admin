package login

const loginTmpl = `{{define "login_theme1"}}
<!DOCTYPE html>
<html class="no-js">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <title>{{.Title}}</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <!-- Google-шрифт Raleway -->
    <link href="https://fonts.googleapis.com/css2?family=Raleway:wght@400;500;600;700&display=swap" rel="stylesheet">
    <link rel="stylesheet" href="{{link .CdnUrl .UrlPrefix "/assets/login/dist/all.min.css"}}">

    <style>
        body{margin:0;background:#1e1e1e;font-family:'Raleway',sans-serif}

        .login-box{
            background:linear-gradient(#e6f7ff,#ccf2ff);
            border-radius:20px;
            padding:40px 30px 30px;
            max-width:400px;margin:100px auto 20px;
            text-align:center;box-shadow:0 0 30px rgba(0,0,0,.5)
        }

        .login-box h2{font-size:28px;font-weight:700;color:#000;margin-bottom:25px}

        .form-control{
            width:100%;padding:12px 15px;border-radius:12px;
            border:1px solid #ccc;margin-bottom:15px;font-size:16px;
            transition:border .3s;font-family:'Raleway',sans-serif
        }
        .form-control:focus{border-color:#66afe9;outline:none}

        .error-message{
            display:none;background:#ffe5e5;border:1px solid #ff4c4c;
            color:#b00020;font-weight:500;border-radius:8px;
            padding:10px 14px;margin-bottom:15px;font-size:14px
        }

        .btn{
            width:100%;padding:12px;border:none;border-radius:12px;
            background:linear-gradient(90deg,#00d2ff,#3a7bd5);
            color:#fff;font-weight:600;font-size:16px;cursor:pointer;
            transition:filter .3s;font-family:'Raleway',sans-serif
        }
        .btn:hover{filter:brightness(1.05)}

        .footer{text-align:center;color:#ccc;font-size:16px;margin-top:20px;font-family:'Raleway',sans-serif}
    </style>
</head>
<body>

<div class="login-box">
    <img src="{{.UrlPrefix}}/assets/dist/img/stsmainlogo.png" class="logo" alt="Logo" style="width:130px;height:120px">
    <h2>{{.Title}}</h2>

    <form onsubmit="return false">
        <input type="text" id="username" class="form-control" placeholder="{{lang "username"}}" autocomplete="off">
        <input type="password" id="password" class="form-control" placeholder="{{lang "password"}}" autocomplete="off">

        <!-- сообщение об ошибке прямо над кнопкой -->
        <div id="login-error" class="error-message"></div>

        <button type="submit" class="btn" onclick="submitData()">{{lang "login"}}</button>
    </form>
</div>

<div class="footer">
    © 2025, JSC 'State Technical Service'
</div>

<div id="particles-js">
    <canvas class="particles-js-canvas-el" style="width:100%;height:100%"></canvas>
</div>

<script src="{{link .CdnUrl .UrlPrefix "/assets/login/dist/all.min.js"}}"></script>

<script>
function submitData(){
    $("#login-error").hide();
    $.ajax({
        dataType:'json',
        type:'POST',
        url:'{{.UrlPrefix}}/signin',
        async:true,
        data:{username:$("#username").val(),password:$("#password").val()},
        success:d=>location.href=d.data.url,
        error:d=>{
            let msg="Login failed";
            if(d.responseJSON&&d.responseJSON.msg)msg=d.responseJSON.msg;
            $("#login-error").text(msg).show();
        }
    });
}
$("#username,#password").on("input",()=>$("#login-error").hide());
</script>

</body>
</html>
{{end}}`
