
<script>
    const roundsPerDay = 5;
    const api_url = 'http://localhost:8080/rand_house';

    // initial state
    let round = 1;
    let imgIdx = 0;
    let welcomePage = true;

    let houseImgs  = [];


    async function fetchHouseData() {
        console.log('Fetching house data');
        try {
            const response = await fetch(api_url);
            let houseData = await response.json();
            console.log(houseData);

            houseImgs = houseData.Imgs;
            houseImgs = houseImgs.filter(img => img !== '');

        } catch (error) {
            console.error('Error fetching data:', error);
        }
    }


    // transition to the game page and fetch house data to start the game
    function handleWelcomeButton() {
        welcomePage = false;
        fetchHouseData();
    }

</script>

<h1>Realista</h1>
<p>O jogo em que adivinhas o preço de casas que não podes comprar</p>


{#if welcomePage}
        <button class="centered-button"  on:click={handleWelcomeButton}>Começar</button>
{:else}
    <div class="container">
        <div class="content-wrapper" >
        <h4>Ronda {round} de {roundsPerDay}</h4>
        <img src={houseImgs[imgIdx]} alt="Description of the image" style="max-width: 100%; max-height: 500px;"/>
    </div>
    </div>
{/if}
<style>
    .container {
        display: flex;
        justify-content: center;
        height: 100%;
        width: 100%;
        border: rgb(0, 0, 0) 1px solid;
    }

    .content-wrapper {
    text-align: center; /* Center aligns the content */
}

    h1 {
      font-size: 2rem;
      text-align: center;
    }

    p {
      font-size: 0.9rem;
      text-align: center;
    }

    .centered-button {
        margin-top: 20px;        /* Adds space between the button and other elements */
        padding: 10px 20px;
        font-size: 16px;
        cursor: pointer;
        display: block;
        margin-left: auto;
        margin-right: auto;
    }
  

  </style>