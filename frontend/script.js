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

// called when the nextHouseButton is clicked.
function nextHouse(){
    console.log("NEXT\n");
    const nextHouseButton = document.getElementById('next-house-button');
    nextHouseButton.style.display = 'none';  // Show the button
    currentImageIndex = 0;
    fetchHouseData();
    guess_counter = 0;
}


// Updates the UI when the page is loaded
function updateUI(data) {
    // get HTML elements
    const areaElement = document.getElementById('area');
    const roomsElement = document.getElementById('rooms');
    const neighborhoodElement = document.getElementById('neighborhood');
    const currentImageElement = document.getElementById('current-image');
    // debug message

    imgs = data.Imgs.filter(function(element){return element != "";}) // the img links vec comes with some garbage
    
    currentImageIndex = 0;
    areaElement.innerText = ` ${data.Area} m²`;
    actualPrice = data.Price / 1000;
    house_link = data.Url;

    // Set the first image in the current-image element
    currentImageElement.src = imgs[currentImageIndex];

    // clear all guess bars to prepare for the new round
    clearGuessBars();
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
        addGuessBar(userGuess,'Correto');
        resultElement.innerHTML = `Spot on!!  <a href="${house_link}"> Check the house here!</a>`;
        roundOver();
        won_games++

    } else if (delta < 20){
        addGuessBar(userGuess,'Correto');
        resultElement.innerHTML = `Close enough! The actual price is ${actualPrice} <a href="${house_link}"> Check the house here!</a>`;
        roundOver();
        won_games++
    } else if(guess_counter >= max_guesses -1){
        addGuessBar(userGuess,'Too many errors');
        resultElement.innerHTML = `You missed too many times! The actual price was ${actualPrice} <a href="${house_link}"> Check the house here!</a>`;
        const nextHouseButton = document.getElementById('next-house-button');
        nextHouseButton.style.display = 'block';  // Show the button
        roundOver();
    } else if(guess_counter == 0){
        if(delta < (actualPrice / 5)) { // tune this value!
            resultElement.innerText = 'Warm'; 
            addGuessBar(userGuess,'Quente');          
        } else {
            resultElement.innerText = 'Cold';
            addGuessBar(userGuess,'Frio')
        }
    } else {
        if(prev_delta < delta) {
            resultElement.innerText = 'Colder...';
            addGuessBar(userGuess,'Mais frio');

        }
        else if(prev_delta > delta) {
            resultElement.innerText = 'Warmer...';
            addGuessBar(userGuess,'Mais quente');
        }
        else {
            resultElement.innerText = 'Change the guess!!';
            return
        }
    }
    guess_counter++;
    prev_delta = delta;
}


function roundOver(){
    const nextHouseButton = document.getElementById('next-house-button');
    nextHouseButton.style.display = 'block';  // Show the button
    total_games++;
    guess_counter = 0;

}

//  to create and add a new guess bar
function addGuessBar(userGuess, hintString) {
    const guessBarsContainer = document.getElementById('guess-bars-container');

    // Create new elements for the guess bar
    const guessBar = document.createElement('div');
    guessBar.classList.add('guess-bar');

    const userGuessElement = document.createElement('p');
    userGuessElement.textContent = `${userGuess}`;

    const hintElement = document.createElement('p');
    hintElement.textContent = hintString;

    // Append elements to the guess bar
    guessBar.appendChild(userGuessElement);
    guessBar.appendChild(hintElement);

    // Append the guess bar to the container
    guessBarsContainer.appendChild(guessBar);
}

// clear all guess bars
function clearGuessBars() {
    const guessBarsContainer = document.getElementById('guess-bars-container');
    guessBarsContainer.innerHTML = ''; // Remove all content inside the container
}