import { Header } from './header.mjs';
import { View, VIEWS } from './view.mjs';
import { Auth } from './auth.mjs';
import { WSClient } from './ws.mjs';

export class Forum {
	constructor() {
		this.header = new Header(this);
		this.view = new View(this);
		this.auth = new Auth(this);
		// after "reload page" button pressed, when the user is logged in redirect to profile page, otherwise redirect/reload to login page
		const profileOrLogin = this.userIsLoggedIn() ? VIEWS.PROFILE : VIEWS.LOGIN;

		this.state = {
			page: profileOrLogin,
			auth: { UUID: null, username: null },
		};
		this.ws = new WSClient(this); // initialize websocket client
		this.ws.socket.onopen = () => this.auth.checkCookie(); // check cookie when websocket connection is established
	}

	userIsLoggedIn() {
		const UUID = (document.cookie.match(/^(?:.*;)?\s*UUID\s*=\s*([^;]+)(?:.*)?$/) || [, null])[1];
		const username = (document.cookie.match(/^(?:.*;)?\s*username\s*=\s*([^;]+)(?:.*)?$/) || [, null])[1];
		if (UUID === null || username === null) {
			return false;
		}
		return true;
	}

	update() {
		this.header.update();
		this.view.update();
	}

	async postSubmit() {
		const data = new FormData(document.getElementById('view-newpost')); // get data from the html form
		fetch('/api/post/submit', { // send data to the server
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ // convert data to JSON for request body
				UUID: this.state.auth.UUID,
				title: data.get('newpost-title'),
				categories: data.get('newpost-categories'),
				text: data.get('newpost-text'),
			}),
		})
			.then((response) => response.json()) // wait for response from server, then convert to JSON...
			.then((data) => { // ...with name "data"
				if (data.message == 'Post created') { // if the server responds with "Post created", then clear the form fields
					document.getElementById('newpost-title').value = '';
					document.getElementById('newpost-categories').value = '';
					document.getElementById('newpost-text').value = '';
					return;
				}
				alert(data.message);
			})
			.catch((err) => {
				console.log(err);
			});
	}

	async commentSubmit() {
		const data = new FormData(document.getElementById('comment-form')); // get data from the html form
		fetch('/api/comment/submit', { // send data to the server
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ // convert data to JSON for request body
				UUID: this.state.auth.UUID,
				PostID: parseInt(this.view.viewPost.postID),
				text: document.getElementById('comment-text').value,
			}),
		})
			.then((response) => response.json()) // wait for response from server, then convert to JSON...
			.then((data) => { // ...with name "data"
				if (data.message == 'Comment created') { // if the server responds with "Comment created", then clear the form field
					document.getElementById('comment-text').value = '';
				} else {
					alert(data.message);
				}
			})
			.catch((err) => {
				console.log(err);
			});
	}

	async messageSubmit() {
		const data = new FormData(document.getElementById('comment-form')); // get data from the html form
		fetch('/api/chat/newmessage', { // send data to the server
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ // convert data to JSON for request body
				usernameFrom: this.state.auth.username,
				usernameTo: this.view.viewMessages.username,
				text: document.getElementById('message-text').value,
			}),
		})
			.then((response) => response.json()) // wait for response from server, then convert to JSON...
			.then((data) => { // ...with name "data"
				if (data.message == 'Message sent') { // if the server responds with "Message sent", then clear the form field
					document.getElementById('message-text').value = '';
				} else {
					alert(data.message);
				}
			})
			.catch((err) => {
				console.log(err);
			});
	}
}
