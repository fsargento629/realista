// Global variables
let actualPrice;  // in kEur
let house_link;
let guess_counter,prev_delta;
let won_games,total_games; 
const max_guesses = 6;
let currentImageIndex = 0;
let imgs;

document.addEventListener("DOMContentLoaded", function () {
    fetchHouseData();
    guess_counter = 0;

    // Event listener for the "Enter" key in the price-input field
    const priceInput = document.getElementById('price-input');
    priceInput.addEventListener('keyup', function (event) {
        if (event.key === 'Enter') {
            checkPrice();
        }
    });
});

function fetchHouseData() {
    // Fetch data from your Go API
    // const api_url = 'http://192.168.1.66:8080/rand_house';
    const api_url = 'http://localhost:8080/rand_house';
 
    fetch(api_url)
        .then(response => response.json())
        .then(data => updateUI(data))
        .catch(error => console.error('Error fetching data:', error));
}

// Updates the UI when the page is loaded
function updateUI(data) {
    // get HTML elements
    const areaElement = document.getElementById('area');
    const roomsElement = document.getElementById('rooms');
    const neighborhoodElement = document.getElementById('neighborhood');
    const currentImageElement = document.getElementById('current-image');

    imgs = data.Imgs.filter(function(element){return element != "";}) // the img links vec comes with some garbage
    
    currentImageIndex = 0;
    areaElement.innerText = ` ${data.Area} m²`;
    roomsElement.innerText = `${data.Rooms}`;
    neighborhoodElement.innerText = data.Bairro.charAt(0).toUpperCase() + data.Bairro.slice(1);
    actualPrice = data.Price / 1000;
    house_link = data.Url;

    // Set the first image in the current-image element
    currentImageElement.src = imgs[currentImageIndex];
}

function showNextImage() {
    currentImageIndex = (currentImageIndex + 1) % imgs.length;
    document.getElementById('current-image').src = imgs[currentImageIndex];
}

function showPreviousImage() {
    currentImageIndex = (currentImageIndex - 1 + imgs.length) % imgs.length;
    document.getElementById('current-image').src = imgs[currentImageIndex];
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
        resultElement.innerHTML = `Spot on!!  <a href="${house_link}"> Check the house here!</a>`;
    } else if (delta < 20){
        resultElement.innerHTML = `Close enough! The actual price is ${actualPrice} <a href="${house_link}"> Check the house here!</a>`;
    } else if(guess_counter >= max_guesses -1){
        resultElement.innerHTML = `You missed too many times! The actual price was ${actualPrice} <a href="${house_link}"> Check the house here!</a>`;
    } else if(guess_counter == 0){
        if(delta < (actualPrice / 5)) { // tune this value!
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
