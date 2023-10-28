import { h, render } from "https://unpkg.com/preact?module";
import htm from "https://unpkg.com/htm?module";

const html = htm.bind(h);

function App(props) {
  return html`
    <div>
    	<p>Hi from Preact</p>
    </div>
  `;
}

document.addEventListener("DOMContentLoaded", () => {
	document.write("<h2>Hello World!</h2><p>Have a nice day!</p>");
	render(html`<${App}>`, document.body);
})

