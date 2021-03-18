class CardData {
    constructor(title, description, imageURLs){
        this.title = title;
        this.description = description;
        this.imageURLs = imageURLs;
    }
}

class Card {
    constructor(data) {
        this.data = data; 
    }

    dom() {
        let imageDiv = document.createElement("DIV");
        imageDiv.classList.add("image");

        let imageURL = this.data.imageURLs[0];
        imageDiv.style.backgroundImage = 'url("' + imageURL + '")';
        imageDiv.style.backgroundRepeat = 'no-repeat';
        imageDiv.style.backgroundSize =  'cover';

        let titleSpan = document.createElement("SPAN");
        titleSpan.classList.add("title");
        titleSpan.appendChild(document.createTextNode(this.data.title));


        let descText = document.createTextNode(this.data.description);
        let newLine = document.createElement("BR");
        let textP = document.createElement("P");
        textP.appendChild(titleSpan);
        textP.appendChild(newLine);
        textP.appendChild(descText);
        

        let textDiv = document.createElement("DIV");
        textDiv.classList.add("text");
        textDiv.appendChild(textP);

        let actionsDiv = document.createElement("DIV");
        actionsDiv.classList.add("actions");
    

        let cardDiv = document.createElement("DIV");
        cardDiv.classList.add("card");
        cardDiv.appendChild(imageDiv);
        cardDiv.appendChild(textDiv);
        cardDiv.appendChild(actionsDiv);
    
        return cardDiv;
    }
}