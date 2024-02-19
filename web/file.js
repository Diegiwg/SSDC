const input = document.querySelector("#distanceInput");
const button = document.querySelector("#calculateButton");

button.addEventListener("click", () => {
  const distance = input.value;
  if (isNaN(distance) || distance <= 0) {
    alert("Please enter a distance greater than 0.");
  }

  document.location.href = "/" + distance;
});
