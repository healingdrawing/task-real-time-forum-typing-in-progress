import { VIEWS } from './view.mjs';

export class Auth {
	constructor(forum) {
		this.forum = forum;
	}

	signup = async function () {
		const data = new FormData(document.getElementById('view-signup'));
		try {
			const response = await fetch('/api/user/register', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					username: data.get('signup-user'),
					age: parseInt(data.get('signup-age')),
					gender: data.get('signup-gender'),
					firstname: data.get('signup-firstname'),
					lastname: data.get('signup-lastname'),
					email: data.get('signup-email'),
					password: data.get('signup-password'),
				}),
			});
			const responseData = await response.json();
			if (Object.hasOwn(responseData, 'UUID')) {
				document.cookie = `username=${responseData.username}`;
				document.cookie = `UUID=${responseData.UUID}`;
				this.forum.state.auth = { UUID: responseData.UUID, username: responseData.username };
				this.forum.state.page = VIEWS.PROFILE;
				this.forum.update();
				this.forum.ws.login();
				return;
			}
			alert(responseData.message);
		} catch (err) {
			console.log(err);
		}
	};

	login = async () => {
		const data = new FormData(document.getElementById('view-login'));
		// fetch data from golang backend
		// convert data from the html form to JSON for request body
		fetch('/api/user/login', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({
				username: data.get('login-user'),
				password: data.get('login-password'),
			}),
		})
			.then((response) => response.json()) // wait for response from backend, and convert to JSON...
			.then((data) => { //  ...with name 'data'
				if (Object.hasOwn(data, 'UUID')) { // if the response has a UUID property, fill the cookies and state with the data
					document.cookie = `username=${data.username}`;
					document.cookie = `UUID=${data.UUID}`;
					this.forum.state.auth = { UUID: data.UUID, username: data.username };
					this.forum.state.page = VIEWS.PROFILE; // switch to profile page after login
					this.forum.update(); // update header and view sections of the page
					this.forum.ws.login(); // login using WSClient to websocket server, from ws.mjs
				}
				alert(data?.username ? `👋 Welcome ${data.username}` : `😕 That login does not exist.`);
			})
			.catch((err) => {
				console.log(err);
			});
	};

	logout = async () => {
		fetch('/api/user/logout', { // fetch data from golang backend
			method: 'DELETE', // method DELETE is used to delete session from backend
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ // convert data from the state to JSON for request body
				UUID: this.forum.state.auth.UUID,
			}),
		})
			.then((response) => response.json()) // wait for response from backend, and convert to JSON...
			.then((data) => { //  ...with name 'data'
				if ((data.message = 'Session deleted')) { // if the message is 'Session deleted', expire cookies
					document.cookie = 'UUID=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
					document.cookie = 'username=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
					this.forum.ws.logout();  // logout using WSClient from websocket server, from ws.mjs
					this.forum.state.auth = { UUID: null, username: null }; // reset state
					this.forum.state.page = VIEWS.LOGIN; // switch to login page after logout
					this.forum.update(); // update header and view sections of the page
					
					// force refresh data with reload the page in logout process
					// Generate a random number or timestamp
					var cacheBuster = new Date().getTime(); // Using timestamp as cache buster

					// Get the current URL
					var url = window.location.href;

					// Append the cache buster query parameter to the URL
					var updatedUrl = url + (url.indexOf('?') === -1 ? '?' : '&') + 'cache=' + cacheBuster;

					// Redirect the browser to the updated URL
					window.location.href = updatedUrl;

					// Clean the URL by removing the cache query parameter
					var cleanedUrl = location.href.replace(/[?&]cache=\d+/, '');

					// Replace the browser URL with the cleaned URL
					history.replaceState(null, null, cleanedUrl);

				}
				alert(data.message);
			})
			.catch((err) => {
				console.log(err);
			});
	};

	checkCookie = async () => {
		const UUID = (document.cookie.match(/^(?:.*;)?\s*UUID\s*=\s*([^;]+)(?:.*)?$/) || [, null])[1];
		const username = (document.cookie.match(/^(?:.*;)?\s*username\s*=\s*([^;]+)(?:.*)?$/) || [, null])[1];
		if (UUID === null || username === null) { // if no matching cookies, reset state
			this.forum.state.auth = { UUID: null, username: null };
			this.forum.update(); // update header and view sections of the page
			return;
		}
		return fetch('/api/user/check', { // fetch data from golang backend
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ // convert data from the cookies to JSON for request body
				UUID,
			}),
		})
			.then((response) => {
				if (response.ok) { // if session exists in database, then response will be includes Exists=true. (backend user.go sessionCheckHandler)
					return response.json(); // wait for response from backend, and convert to JSON...
				}
				throw response.statusText;
			})
			.then((data) => { //  ...with name 'data'
				if (data.Exists) {
					this.forum.state.auth = { UUID: UUID, username: username };
					this.forum.update();
					this.forum.ws.login();
					return;
				}
				document.cookie = 'UUID=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
				document.cookie = 'username=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
				this.forum.state.auth = { UUID: null, username: null };
				this.forum.update();
			})
			.catch((err) => {
				alert(err);
				console.log(err);
			});
	};
}
