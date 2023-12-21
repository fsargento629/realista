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
    document.getElementById('rooms').innerText = data.Rooms;
    document.getElementById('neighborhood').innerText = data.Bairro;
    document.getElementById('actual-price').innerText = data.Price;
}

function checkPrice() {
    // Get user's guess
    const userGuess = parseInt(document.getElementById('price-input').value);

    // Get the actual price from the displayed UI
    const actualPrice = parseInt(document.getElementById('actual-price').innerText);

    // Compare user's guess with actual price
    const resultElement = document.getElementById('result');
    if (userGuess === actualPrice) {
        resultElement.innerText = 'Congratulations! Your guess is correct!';
    } else {
        resultElement.innerText = `Oops! Your guess is incorrect. The actual price is ${actualPrice}.`;
    }
}
