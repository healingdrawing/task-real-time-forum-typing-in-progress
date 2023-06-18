import { checkSignup, checkLogin, checkPost, checkComment, checkMessage } from './forms.mjs';
import { throttle } from './utilities.mjs';

export const VIEWS = {
	SIGNUP: 0,
	LOGIN: 1,
	POSTS: 2,
	POST: 3,
	CHAT: 5,
	PROFILE: 6,
};

export class View {
	constructor(forum) {
		this.forum = forum;

		this.comments = document.getElementById('comments');
		this.messageBank = [];
		this.messages = document.getElementById('messages');
		this.messageTyping = document.getElementById('message-typing');
		this.posts = document.getElementById('posts');
		this.unread = [];
		this.users = [];
		this.viewChat = document.getElementById('view-chat');
		this.viewMessages = document.getElementById('view-messages');
		this.viewLogin = document.getElementById('view-login');
		this.viewPost = document.getElementById('view-post');
		this.viewPosts = document.getElementById('view-posts');
		this.viewProfile = document.getElementById('view-profile');
		this.viewSignup = document.getElementById('view-signup');
		this.viewUsers = document.getElementById('view-users');

		document.getElementById('view-signup').addEventListener('submit', (e) => {
			if (checkSignup()) this.forum.auth.signup();
			e.preventDefault();
		});

		document.getElementById('view-login').addEventListener('submit', (e) => {
			if (checkLogin()) this.forum.auth.login();
			e.preventDefault();
		});

		document.getElementById('view-newpost').addEventListener('submit', (e) => {
			if (checkPost()) this.forum.postSubmit();
			e.preventDefault();
		});

		document.getElementById('comment-form').addEventListener('submit', (e) => {
			if (checkComment()) this.forum.commentSubmit();
			e.preventDefault();
		});

		document.getElementById('message-form').addEventListener('submit', (e) => {
			if (checkMessage()) this.forum.messageSubmit();
			e.preventDefault();
		});

		this.typing = false;
		document.getElementById('message-form').addEventListener('input', (e) => {
			if (!this.typing) {
				fetch("/api/chat/typing", {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({
						usernameFrom: this.forum.state.auth.username,
						usernameTo: this.viewMessages.username,
						typing: true,
					}),
				})
					.then((response) => response.json())
					.then((data) => {
						if (data.message != "OK") alert(data.message);
					})
					.catch((err) => {
						console.log(err);
					});
			}
			else clearTimeout(this.typing);
			this.typing = setTimeout(() => {
				fetch("/api/chat/typing", {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({
						usernameFrom: this.forum.state.auth.username,
						usernameTo: this.viewMessages.username,
						typing: false,
					}),
				})
					.then((response) => response.json())
					.then((data) => {
						if (data.message != "OK") alert(data.message);
						else this.typing = false;
					})
					.catch((err) => {
						console.log(err);
					});
			}, 500);
		});

		document.getElementById('messages').addEventListener(
			'scroll',
			throttle(
				() => {
					if (document.getElementById('messages').scrollTop == 0) this.messagesLoadMore();
				},
				500,
				{ leading: true, trailing: true }
			)
		);
	}

	update() {
		this.hideAll();
		switch (this.forum.state.page) {
			case VIEWS.SIGNUP:
				document.getElementById('signup-user').value = '';
				document.getElementById('signup-age').value = '';
				document.getElementById('signup-gender').value = 'undefined';
				document.getElementById('signup-firstname').value = '';
				document.getElementById('signup-lastname').value = '';
				document.getElementById('signup-email').value = '';
				document.getElementById('signup-password').value = '';
				this.viewSignup.hidden = false;
				this.viewUsers.hidden = true;
				break;
			case VIEWS.LOGIN:
				document.getElementById('login-user').value = '';
				document.getElementById('login-password').value = '';
				this.viewLogin.hidden = false;
				this.viewUsers.hidden = true;
				this.viewChat.hidden = true;
				this.viewMessages.hidden = true;
				this.viewProfile.hidden = true;
				this.viewSignup.hidden = true;
				this.viewPosts.hidden = true;
				this.viewPost.hidden = true;
				break;
			case VIEWS.POSTS:
				document.getElementById('newpost-title').value = '';
				document.getElementById('newpost-categories').value = '';
				document.getElementById('newpost-text').value = '';
				document.getElementById('view-newpost').hidden = this.forum.state.auth.UUID ? false : true;
				this.fetchPosts();
				this.viewPosts.hidden = false;
				this.viewUsers.hidden = false;
				this.fetchUsers();
				break;
			case VIEWS.POST:
				document.getElementById('comment-text').value = '';
				document.getElementById('comment-form').hidden = this.forum.state.auth.UUID ? false : true;
				this.fetchComments();
				this.viewPost.hidden = false;
				this.viewUsers.hidden = false;
				this.fetchUsers();
				break;
			case VIEWS.CHAT:
				this.fetchUsers();
				this.viewMessages.hidden = true;
				this.viewUsers.hidden = false;
				this.viewChat.hidden = false;
				break;
			case VIEWS.PROFILE:
				this.viewProfile.hidden = false;
				this.viewUsers.hidden = false;
				this.fetchUsers();
				this.showProfile(this.forum.state.auth.username);
				break;
		}
	}

	showChat() {
		this.viewMessages.hidden = true;
		this.viewUsers.hidden = false;
		this.viewChat.hidden = false;
	}

	fetchPosts() {
		fetch('/api/post/get', {
			method: 'GET',
		})
			.then((response) => response.json())
			.then((data) => {
				if (!Object.hasOwn(data, 'posts')) {
					alert(data.message);
				} else if (data.posts != null) {
					this.posts.replaceChildren(...data.posts.map((post) => this.generatePost(post)));
				} else {
					document.querySelectorAll('.comment').forEach((e) => e.remove());
				}
			})
			.catch((err) => {
				console.log(err);
			});
	}

	fetchComments() {
		fetch('/api/comment/get', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({
				postID: parseInt(this.viewPost.postID),
			}),
		})
			.then((response) => response.json())
			.then((data) => {
				if (!Object.hasOwn(data, 'comments')) {
					alert(data.message);
				} else if (data.comments != null) {
					this.comments.replaceChildren(
						...data.comments.map((comment) => this.generateComment(comment))
					);
				} else {
					document.querySelectorAll('.comment').forEach((e) => e.remove());
				}
			})
			.catch((err) => {
				console.log(err);
			});
	}

	showProfile(username) {
		fetch('/api/user/profile', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({
				username,
			}),
		})
			.then((response) => response.json())
			.then((data) => {
				if (Object.hasOwn(data, 'message')) {
					alert(data.message);
					return;
				}
				this.forum.state.page = VIEWS.PROFILE;
				this.hideAll();
				this.viewProfile.hidden = false;
				this.viewUsers.hidden = false;
				this.viewProfile.innerText = `Username: ${data.username}
					Age: ${data.age}
					Gender: ${data.gender}
					First name: ${data.firstname}
					Last name: ${data.lastname}
					E-mail: ${data.email}`;
			})
			.catch((err) => {
				console.log(err);
			});
	}

	hideAll() {
		this.viewSignup.hidden = true;
		this.viewLogin.hidden = true;
		this.viewPosts.hidden = true;
		this.viewPost.hidden = true;
		this.viewChat.hidden = true;
		this.viewProfile.hidden = true;
		this.viewUsers.hidden = true;
	}

	fetchUsers() {
		fetch('/api/chat/getusers', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({
				username: this.forum.state.auth.username,
			}),
		})
			.then((response) => response.json())
			.then((data) => {
				if (Object.hasOwn(data, 'users')) {
					this.users = Object.fromEntries(
						data.users
							.filter((user) => user.username != this.forum.state.auth.username)
							.map((user) => [user.username, { data: user, element: this.generateUser(user), typing: false }])
					);
					this.updateUserList();
					this.updateNotifications();
					return;
				}
				alert(data.message);
			})
			.catch((err) => {
				console.log(err);
			});
	}

	generatePost = (data) => {
		const post = document.createElement('div');
		const user = document.createElement('div');
		const title = document.createElement('div');
		const categories = document.createElement('div');
		const text = document.createElement('div');
		post.className = 'post';
		user.className = 'post-username';
		title.className = 'post-title';
		categories.className = 'post-categories';
		user.innerText = data.username;
		title.innerText = data.title;

		data.categories.split(',').forEach((cat) => {
			const category = document.createElement('div');
			category.className = 'post-category';
			category.innerText = cat;
			categories.appendChild(category);
		});
		text.innerText = data.text;
		text.hidden = true;
		post.appendChild(user);
		post.appendChild(title);
		post.appendChild(categories);
		post.appendChild(text);
		post.onclick = () => {
			this.viewPost.postID = data.ID;
			document.getElementById('view-post-post').replaceChildren(post);
			text.hidden = false;
			this.forum.state.page = VIEWS.POST;
			this.update();
		};
		return post;
	};

	generateComment = (data) => {
		const comment = document.createElement('div');
		const user = document.createElement('div');
		const text = document.createElement('div');
		comment.className = 'comment';
		user.className = 'comment-user';
		text.className = 'comment-text';
		user.innerText = data.username;
		text.innerText = data.text;
		comment.appendChild(user);
		comment.appendChild(text);
		return comment;
	};

	generateUser = (data) => {
		const user = document.createElement('div');
		const userProfile = document.createElement('div');
		const container = document.createElement('div');
		userProfile.innerText = `${data.online ? '😀' : '😴'}`;
		userProfile.onclick = () => {
			this.showProfile(data.username);
		};
		userProfile.className = 'profile';
		user.id = `user-${data.username}`;
		user.className = 'user';
		container.innerText = `${data.username}`;
		if (data.time && new Date(data.time) > 8) {
			const time = document.createElement('div');
			time.className = 'message-from-time';
			time.innerText = new Date(data.time).toLocaleString();
			container.appendChild(time);
		}
		user.appendChild(userProfile);
		user.appendChild(container);
		user.onclick = () => {
			this.messageBank = [];
			this.showChat();
			this.viewMessages.username = data.username;
			this.viewMessages.hidden = false;
			document.getElementById('message-text').innerText = '';
			this.messages.innerText = `There are no messages between you and "${data.username}" yet.`;
			fetch('/api/chat/getmessages', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					username: this.forum.state.auth.username,
					otheruser: data.username,
				}),
			})
				.then((response) => response.json())
				.then((responseData) => {
					if (!Object.hasOwn(responseData, 'messages')) alert(responseData.message);
					else if (responseData.messages) {
						this.messageBank = responseData.messages
							.reverse()
							.map((message) => this.generateMessage(message));
						if (this.messageBank.length) {
							this.messages.innerText = '';
							this.messagesLoadMore();
						}
						this.messages.scrollTop = this.messages.scrollHeight;
					}
					if (this.unread.includes(data.username)) {
						this.unread = this.unread.filter((username) => username != data.username);
						this.updateNotifications();
					}
					this.updateTyping()
				})
				.catch((err) => {
					console.log(err);
				});
		};
		return user;
	};

	generateMessage = (data) => {
		const message = document.createElement('div');
		message.className = 'message';
		const messageFrom = document.createElement('div');
		messageFrom.className = 'message-from';
		const username = document.createElement('div');
		username.className = 'message-from-user';
		const time = document.createElement('div');
		time.className = 'message-from-time';
		const messageText = document.createElement('div');
		messageText.className = 'message-text';
		username.innerText = data.usernameFrom;
		time.innerText = new Date(data.time).toLocaleString();
		messageText.innerText = data.text;
		messageFrom.appendChild(username);
		messageFrom.appendChild(time);
		message.appendChild(messageFrom);
		message.appendChild(messageText);
		return message;
	};

	updateUserList = () => {
		const usersWithTime = [];
		const usersWithoutTime = [];

		Object.entries(this.users).forEach((user) => {
			if (user[1]?.data?.time == null) {
				usersWithoutTime.push(user);
			}
			if (user[1]?.data?.time == "0001-01-01T00:00:00Z") {
				usersWithoutTime.push(user);
			} else {
				usersWithTime.push(user);
			}
		});

		usersWithTime.sort((u1, u2) => u1[0].localeCompare(u2[0]));
		usersWithoutTime.sort((u1, u2) => u1[0].localeCompare(u2[0]));

		usersWithTime.sort((u1, u2) => {
			if (!u1[1].data.time) return 1;
			if (!u2[1].data.time) return -1;
			return new Date(u2[1].data.time).getTime() - new Date(u1[1].data.time).getTime();
		});

		const combinedUsers = [...usersWithTime, ...usersWithoutTime];

		if ((combinedUsers.length) != (this.viewUsers.childElementCount)) {
			for (let i = 0; i < combinedUsers.length; i++) {
				if (!this.viewUsers.children[i]) {
					this.viewUsers.appendChild(combinedUsers[i][1].element);
				}
			}
		}

		this.viewUsers.replaceChildren(...combinedUsers.map((user) => user[1].element));
	};


	updateNotifications() {
		this.forum.header.linkChat.innerText = this.unread.length ? '📬 to chat' : '📭 to chat';
		if (this.users)
			Object.entries(this.users).forEach(([username, data]) => {
				if (this.unread.includes(username)) data.element.classList.add('chat-unread');
				else data.element.classList.remove('chat-unread');
			});
		this.updateUserList()
	}

	messagesLoadMore() {
		if (this.messages.children.length < this.messageBank.length)
			this.messageBank
				.slice(
					Math.max(this.messageBank.length - this.messages.children.length - 10, 0),
					this.messageBank.length - this.messages.children.length
				)
				.reverse()
				.forEach((message) => {
					this.messages.prepend(message);
					this.messages.scrollTop += document.getElementsByClassName('message')[0].clientHeight;
				});
	}

	updateTyping() {
		if (this?.users[this?.viewMessages?.username]?.data?.typing) {
			document.getElementById('message-typing').hidden = false
			document.getElementById('message-typing-text').innerText = `${this.viewMessages.username} is typing`
		}
		else document.getElementById('message-typing').hidden = true
	}
}
