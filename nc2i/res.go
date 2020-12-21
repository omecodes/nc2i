package nc2i

import (
	"bytes"
	ht "html/template"
	tt "text/template"
)

type EmailData struct {
	Year    string `json:"year"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Details string `json:"details"`
}

var messageEmailHTMLContent = `
<html>
<head>
    <meta charset="UTF-8">
    <title>NC2I: Demande de devis</title>
	<style>
		body {
			margin: 0;
			padding: 0;
			font-size: 24px;
			color: #212121;
			background-color: white;
		}
	
		#content {
			width: 350px;
			margin: 0 auto;
		}
	
		.main-colored {
			color: #0D47A1;
		}
	
		.text {
			color: #444444;
		}
	
		.details-box  {
			display: block;
			width: 300px;
			padding: 32px;
			background-color: #EEEEEE;
			border-radius: 4px;
		}
	
		.name {
			font-weight: bold;
		}
	
		.title {
			font-size: 24px; 
		}
	
		.footer {
			width: 100%;
			padding: 32px;
		}
	
		.footer p {
			width: 300px;
			margin: 0 auto;
			font-size: 12px;
		}
	
		.section-space {
			height: 30px;
		}
	</style>
</head>
<body>
<div id="content">
    <p class="title main-colored">Nouvelle demande de devis</p>
    <p class="message-head"><span class="name">{{.Name}}</span> a demandé devis en laissant le message suivant: </p>

    <p class="details-box">{{.Details}}</p>
	</br>

	</br>
    <p class="message-head">Avec les coordonnées suivantes:</p>
	</br>
    <table class="contact">
        <tr>
            <td>Email</td>
            <td><a class="button text" href="mailto: {{.Email}}">{{.Email}}</a></td>
        </tr>
        <tr>
            <td>Tel</td>
            <td><a class="button text" href="tel: {{.Phone}}">{{.Phone}}</a></td>
        </tr>
    </table>
</div>
<div class="section-space"></div>
<div class="footer">
    <p>Omecodes {{.Year}} - Tous droits réservés</p>
</div>
</body>
</html>
`

var messageEmailTextContent = `
    Nouvelle demande de devis

    {{.Name}} a demandé un devis en laissant le message suivant:

    {{.Details}}


    et a laissé les coordonnées suivantes:
	Tel  : {{.Email}}
	Email: {{.Phone}}

    Omecodes {{.Year}} - Tous droits réservés
`

func getEmails(data interface{}) (plain string, html string, error error) {
	hTpl, err := ht.New("message").Parse(messageEmailHTMLContent)
	if err != nil {
		return "", "", err
	}
	hBuff := &bytes.Buffer{}
	err = hTpl.Execute(hBuff, data)
	if err != nil {
		return
	}

	tTpl, err := tt.New("message").Parse(messageEmailTextContent)
	if err != nil {
		return
	}
	tBuff := &bytes.Buffer{}
	err = tTpl.Execute(tBuff, data)
	if err != nil {
		return
	}

	html = string(hBuff.Bytes())
	plain = string(tBuff.Bytes())

	return
}
