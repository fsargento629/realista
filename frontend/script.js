// Global variables
let actualPrice; 

document.addEventListener("DOMContentLoaded", function () {
    fetchHouseData();
});

function fetchHouseData() {
    // Fetch data from your Go API
    fetch('http://localhost:8080/rand_house')
        .then(response => response.json())
        .then(data => updateUI(data))
        .catch(error => console.error('Error fetching data:', error));
}

function updateUI(data) {
    document.getElementById('house-image').src = data.Imgs[1]; // Assuming 'Url' is the property that contains the image URL
    console.log(data.Imgs)
    document.getElementById('area').innerText = data.Area;
    console.log(data.Area)
    document.getElementById('rooms').innerText = data.Rooms;
    console.log(data.Rooms)
    document.getElementById('neighborhood').innerText = data.Bairro;
    console.log(data.Bairro)
    actualPrice = data.Price
    console.log(data.Price)
}

function checkPrice() {
    // Get user's guess
    const userGuess = parseInt(document.getElementById('price-input').value);
    console.log(userGuess)

    // Alternatively, you can directly use the price from the API response
    console.log(actualPrice)

    // Compare user's guess with actual price
    const resultElement = document.getElementById('result');
    if (userGuess === actualPrice) {
        resultElement.innerText = 'Congratulations! Your guess is correct!';
    } else {
        resultElement.innerText = `Oops! Your guess is incorrect. The actual price is ${actualPrice}.`;
    }
}
