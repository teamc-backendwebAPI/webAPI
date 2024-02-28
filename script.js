let target = document.getElementById('target');
const recipeName = 'チキンカレー';
console.log('fetching recipe: ' + target);
fetch('http://localhost:8080/recipe?name=' + recipeName)
  .then((response) => response.json())
  .then((data) => {
    console.log(data);
    target.innerHTML += createRecipeHTML(data);
  });

function createRecipeHTML(data) {
  let ingredientsHtml = data.ingredients
    .map((ingredient) => `<li>${ingredient.name} - ${ingredient.amount}</li>`)
    .join('');

  let stepsHtml = data.steps.map((step, index) => `<li>${step}</li>`).join('');

  let child = `
    <div class="card">
        <div class="card-body">
            <h5 class="card-title">${data.name}</h5>
            <p class="card-text">${data.description}</p>
            <img src=${data.image} class="img-thumbnail" alt="..." />
            <ul class="list-group list-group-flush">
                <li class="list-group-item">Category: ${data.categories.join(
                  ', '
                )}</li>
                <li class="list-group-item">
                    Ingredients:
                    <ul>
                        ${ingredientsHtml}
                    </ul>
                </li>
                <li class="list-group-item">
                    Steps:
                    <ol>
                        ${stepsHtml}
                    </ol>
                </li>
                <li class="list-group-item">
                    Nutrition: Calories ${data.nutrition.calories}, Protein ${
    data.nutrition.protein
  }, Fat ${data.nutrition.fat}, Carbohydrates ${data.nutrition.carbohydrates}
                </li>
                <li class="list-group-item">Difficulty: ${data.difficulty}</li>
                <li class="list-group-item">
                    Prep Time: ${data.time.prep} mins, Cook Time: ${
    data.time.cook
  } mins
                </li>
            </ul>
            <p class="card-text">
                <small class="text-muted">Last updated 3 mins ago</small>
            </p>
        </div>
    </div>
`;

  return child;
}
