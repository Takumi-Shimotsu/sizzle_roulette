// 画像ギャラリー
const mainImage = document.querySelector(".gallery-image img");
const thumbImages = document.querySelectorAll(".gallery-thumbnails img");

thumbImages.forEach((thumbImage) => {
  thumbImage.addEventListener("mouseover", (event) => {
    mainImage.src = event.target.src;
    mainImage.animate({ opacity: [0, 1] }, 5000);
  });
});

// ボタン
// 予算
const budgetButtons = document.querySelectorAll(".budget-buttons button");
let budget = 0; //letは再宣言不可能
budgetButtons.forEach((budgetButton) => {
  budgetButton.addEventListener("click", (event) => {
    //変数に格納
    budget = event.target.value;
    var elm_budget = document.getElementById("budget-text");
    elm_budget.textContent = budget;
    budget = event.target.value;
    if (
      sizzleTaste.length &&
      sizzleTexture.length &&
      sizzleInformation.length &&
      budget > 0
    ) {
      querySend.disabled = false;
    }
  });
});

// しずる
const sizzleButtons = document.querySelectorAll(".grid-buttons button");
const querySend = document.getElementById("start");
let sizzleTaste = [];
let sizzleTexture = [];
let sizzleInformation = [];
sizzleButtons.forEach((sizzleButton) => {
  sizzleButton.addEventListener("click", (event) => {
    if (event.target.className === "tasteButton") {
      sizzleTaste.push(event.target.value);
      var text = document.getElementById("taste-text");
      text.textContent = sizzleTaste;
    }
    if (event.target.className === "textureButton") {
      sizzleTexture.push(event.target.value);
      var text = document.getElementById("texture-text");
      text.textContent = sizzleTexture;
    }
    if (event.target.className === "informationButton") {
      sizzleInformation.push(event.target.value);
      var text = document.getElementById("information-text");
      text.textContent = sizzleInformation;
    }
    sizzleButton.disabled = true;
    if (
      sizzleTaste.length &&
      sizzleTexture.length &&
      sizzleInformation.length &&
      budget > 0
    ) {
      querySend.disabled = false;
    }
  });
});

// しずるリセット
const sizzleReset = document.getElementById("reset");
sizzleReset.addEventListener("click", () => {
  sizzleButtons.forEach((sizzleButton) => {
    sizzleButton.disabled = false;
  });
  querySend.disabled = true;

  var elm_budget = document.getElementById("budget-text");
  elm_budget.textContent = "〇〇";
  budget = 0;

  var elm_sizzle = document.getElementById("taste-text");
  elm_sizzle.textContent = "＜食感系＞";
  var elm_sizzle = document.getElementById("texture-text");
  elm_sizzle.textContent = "＜味覚系＞";
  var elm_sizzle = document.getElementById("information-text");
  elm_sizzle.textContent = "＜情報系＞";
  sizzleTaste = [];
  sizzleTexture = [];
  sizzleInformation = [];
});

// スタート
querySend.addEventListener("click", () => {
  let url = "./result";
  url = url.concat(
    "?budget=",
    budget,
    "&taste=",
    sizzleTaste,
    "&texture=",
    sizzleTexture,
    "&information=",
    sizzleInformation
  );
  // ページ遷移
  window.location.href = url; // 通常の遷移
});
