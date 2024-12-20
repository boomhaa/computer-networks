let button = document.getElementById("sendButton")

button.onclick = () => login()


function login() {
    let host = document.getElementById("host").value
    let login = document.getElementById("login").value
    let password = document.getElementById("password").value
    fetch("/login",{
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            host:host,
            login:login,
            password:password
        },
        )
}).then(response =>  {
    if (response.ok){
        window.location = `http://${window.location.host}/work`
    }else {
        if (response.status===401){
            alert("Invalid credentials")
        }

    }
    }).catch(e =>{
        console.log("Ошибка: " + e)
    })


}