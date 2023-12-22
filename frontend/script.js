// Global variables
let actualPrice;  // in kEur
let guess_counter,prev_delta;
let won_games,total_games; 
const max_guesses = 6;

document.addEventListener("DOMContentLoaded", function () {
    fetchHouseData();
    guess_counter = 0;
});

function fetchHouseData() {
    // Fetch data from your Go API
    fetch('http://localhost:8080/rand_house')
        .then(response => response.json())
        .then(data => updateUI(data))
        .catch(error => console.error('Error fetching data:', error));
}

// Updates the UI when the page is loaded
function updateUI(data) {
    document.getElementById('house-image').src = data.Imgs[1]; // Assuming 'Url' is the property that contains the image URL
    console.log(data.Imgs)
    document.getElementById('area').innerText = `${data.Area} m²`;
    console.log(data.Area)
    document.getElementById('rooms').innerText = data.Rooms;
    console.log(data.Rooms)
    document.getElementById('neighborhood').innerText = data.Bairro.charAt(0).toUpperCase() + data.Bairro.slice(1);
    console.log(data.Bairro)
    actualPrice = data.Price/1000
    console.log(data.Price)
}



function checkPrice() {
    // Get user's guess
    const userGuess = parseInt(document.getElementById('price-input').value);


    // Compare user's guess with actual price
    const resultElement = document.getElementById('result');
    const delta = Math.abs(userGuess - actualPrice);
    console.log(actualPrice); // just for debug, delete this


    // These conditionals are a mess, there must be a better way... REFACTOR
    if ( delta == 0) {
        resultElement.innerText = 'Spot on!';
    } else if (delta < 20){
        resultElement.innerText = `Close enough! The actual price is ${actualPrice}`;
    } else if(guess_counter >= max_guesses -1){
        resultElement.innerText = `You missed too many tries! The actual price was ${actualPrice}`;
    } else if(guess_counter == 0){
        if(delta < 100) {
            resultElement.innerText = 'Warm';           
        } else {
            resultElement.innerText = 'Cold';
        }
    } else {
        if(prev_delta < delta) {
            resultElement.innerText = 'Colder...';
        }
        else if(prev_delta > delta) {
            resultElement.innerText = 'Warmer...';
        }
        else {
            resultElement.innerText = 'Change the guess!!';
            return
        }
    }
    guess_counter++;
    prev_delta = delta;
}
