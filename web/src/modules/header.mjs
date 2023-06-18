import { VIEWS } from './view.mjs';

export class Header {
	constructor(forum) {
		this.forum = forum;
		this.linkChat = document.getElementById('link-chat'); // link to chat page, when you want to see only chat

		document.getElementById('link-home').onclick = () => { // link to forum
			this.forum.state.page = VIEWS.POSTS;
			this.forum.view.update(); // collect data and update the view
		};

		document.onload = () => {
			this.forum.state.page = VIEWS.CHAT;
			this.forum.view.update();
		};

		this.linkChat.onclick = () => {
			this.forum.state.page = VIEWS.CHAT;
			this.forum.view.update();
		};

		document.getElementById('button-logout').onclick = () => {
			this.forum.auth.logout();
		};

		document.getElementById('button-login').onclick = () => {
			this.forum.state.page = VIEWS.LOGIN;
			this.forum.view.update();
		};

		document.getElementById('button-signup').onclick = () => {
			this.forum.state.page = VIEWS.SIGNUP;
			this.forum.view.update();
		};
		document.getElementById('auth-in-user-info').onclick = () => {
			this.forum.view.showProfile(this.forum.state.auth.username)
		};
	}

	update() {
		if (this.forum.state.auth.UUID) {
			document.getElementById('auth-out').hidden = true;
			document.getElementById('auth-in').hidden = false;
			document.getElementById('auth-in-user-info').innerText = `Logged in as ${this.forum.state.auth.username}`;
		} else {
			document.getElementById('auth-out').hidden = false;
			document.getElementById('auth-in').hidden = true;
		}
	}
}
