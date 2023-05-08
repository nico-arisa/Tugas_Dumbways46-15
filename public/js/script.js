
function getData(){
    let name = document.getElementById("name").value
    let email = document.getElementById("email").value
    let phone = document.getElementById("phone").value
    let subject = document.getElementById("subject").value
    let message = document.getElementById("message").value
    
// validation
    if(name == ""|| !name){
        return alert("Nama Harus Di isi!")
    } else if(email == ""){
        return alert("email Harus Di isi!")
    } else if(phone == ""){
        return alert("No telp. Harus Di isi!")
    } else if(subject == ""){
        return alert("subject Harus Di isi!")
    } else if(message == ""){
        return alert("pesan Harus Di isi!")
    }


 
    const destination = "nicoarisaputra@gmail.com"
    let a = document.createElement("a")
    a.href = `mailto:${destination}?subject=${subject}&body= Hallo nama saya ${name} , saya ingin ${message}, bisakah anda menghubungi saya di ${phone}`
    a.click()

    let data = {
        nama: name,
        email: email,
        telp: number,
        subject: subject,
        pesan: message,
    }

    console.log(data)
}