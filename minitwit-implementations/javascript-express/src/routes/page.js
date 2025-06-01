const express = require('express');

const {
    getUserIdByName,
    getUserTimelineMessages,
    getPublicMessages,
    getUser,
    getUserByName,
    isFollowing,
    getUserMessages,
    followUser,
    unfollowUser,
    createNewUser,
    createMessage,
} = require('../database/repository');

const { formatMessages, validateRegisterFields } = require('../utils/utils');

const PER_PAGE = 30;

const router = express.Router();

router.get('/', async (req, res) => {
    try {
        let user = req.session?.user;

        if (!user) {
            return res.redirect('/public');
        }

        const userRefresh = (await getUser(user.user_id)) ?? user;

        const messages = await getUserTimelineMessages(
            userRefresh.user_id,
            PER_PAGE
        );

        res.render('timeline.ejs', {
            user: userRefresh,
            messages,
            title: 'My Timeline',
            flashes: req.flash('success'),
            endpoint: 'myTimeline',
        });
    } catch (error) {
        console.error('Error: ', error);

        res.status(500).send('Internal Server Error');
    }
});

router.get('/public', async (req, res) => {
    try {
        let user = req.session.user;

        const userRefresh = user ? await getUser(user.user_id) : user;

        const messages = formatMessages(await getPublicMessages(PER_PAGE));

        res.render('timeline.ejs', {
            user: userRefresh,
            messages: messages,
            title: 'Public Timeline',
            flashes: req.flash('success'),
            endpoint: 'publicTimeline',
        });
    } catch (error) {
        console.error('Error:', error);

        res.status(500).send('Internal Server Error');
    }
});

router.get('/register', async (req, res) => {
    if (req.session.user) {
        return res.redirect('/timeline');
    }

    res.render('register.ejs', { error: null, flashes: req.flash('success') });
});

router.post('/register', async (req, res) => {
    try {
        const { username, email, password, password2 } = req.body;

        const invalidUserError = await validateRegisterFields(
            username,
            email,
            password,
            password2
        );

        if (invalidUserError) {
            return res.render('register.ejs', {
                error: invalidUserError,
                flashes: req.flash('success'),
            });
        } else if (await getUserIdByName(username)) {
            return res.render('register.ejs', {
                error: 'The username is already taken',
                flashes: req.flash('success'),
            });
        }

        createNewUser(username, email, password);

        req.flash(
            'success',
            'You were successfully registered and can login now'
        );

        res.redirect('/login');
    } catch (error) {
        console.error('Error:', error);

        res.status(500).send('Internal Server Error');
    }
});

router.get('/login', (req, res) => {
    if (req.session.user) {
        return res.redirect('/timeline');
    }

    res.render('login.ejs', {
        user: req.session.user,
        error: null,
        flashes: req.flash('success'),
    });
});

router.post('/login', async (req, res) => {
    try {
        const { username, password } = req.body;

        if (!username || !password) {
            return res.render('login.ejs', {
                error: 'Empty fields not allowed',
                flashes: req.flash('success'),
            });
        }

        const user = await getUserByName(username);

        if (!user) {
            return res.render('login.ejs', {
                error: 'Invalid username',
                flashes: req.flash('success'),
            });
        }

        if (password !== user.pw_hash) {
            return res.render('login.ejs', {
                error: 'Invalid password',
                flashes: req.flash('success'),
            });
        }

        req.session.user = user;

        req.session.save();

        req.flash('success', 'You were logged in');

        res.redirect('/');
    } catch (error) {
        console.error('Error:', error);
        res.status(500).send('Internal Server Error');
    }
});

router.get('/logout', async (req, res) => {
    req.flash('success', 'You were logged out');
    delete req.session.user;
    res.redirect('/public');
});

router.post('/add_message', async (req, res) => {
    const user = req.session.user;

    if (!user) {
        return res.status(401).send('Unauthorized');
    }

    const text = req.body.text;

    if (!text) {
        req.flash('success', 'Message cannot be empty!'); // minimal effort maximum result ¯\_(ツ)_/¯

        return res.redirect('back');
    }

    const success = await createMessage(user.user_id, text);

    if (success) {
        req.flash('success', 'Your message was recorded');

        res.redirect('back');
    } else {
        res.status(500).send('Internal Server Error');
    }
});

router.get('/user/:username', async (req, res) => {
    try {
        let user = req.session?.user;

        const userRefresh = user ? await getUser(user.user_id) : user;

        let profileUser = await getUserByName(req.params.username);

        if (!profileUser) {
            return res.status(404).send('User not found');
        }

        const followed = userRefresh
            ? await isFollowing(userRefresh, profileUser)
            : false;

        const messages = formatMessages(
            await getUserMessages(profileUser.user_id)
        );

        res.render('timeline.ejs', {
            messages: messages,
            followed: followed,
            user: userRefresh,
            profileUser: profileUser,
            title: profileUser.username + "'s Timeline",
            flashes: req.flash('success'),
            endpoint: 'userTimeline',
        });
    } catch (error) {
        console.error('Error fetching user timeline:', error);

        res.status(500).send('Internal Server Error');
    }
});

router.get('/:username/follow', async (req, res) => {
    try {
        const user = req.session.user;

        if (!user) {
            return res.status(401).send('Unauthorized');
        }

        const whomUsername = req.params.username;

        const whomId = await getUserIdByName(whomUsername);

        if (!whomId) {
            return res.status(404).send('User not found');
        }

        if (await followUser(user.user_id, whomId))
            req.flash('success', `You are now following ${whomUsername}`);

        res.redirect(`/user/${whomUsername}`);
    } catch (error) {
        console.error('Error:', error);

        res.status(500).send('Internal Server Error');
    }
});

router.get('/:username/unfollow', async (req, res) => {
    try {
        const user = req.session.user;

        if (!user) {
            return res.status(401).send('Unauthorized');
        }

        const whomUsername = req.params.username;

        const whomId = await getUserIdByName(whomUsername);

        if (!whomId) {
            return res.status(404).send('User not found');
        }

        if (unfollowUser(user.user_id, whomId))
            req.flash('success', `You are no longer following ${whomUsername}`);

        res.redirect(`/user/${whomUsername}`);
    } catch (error) {
        console.error('Error:', error);

        res.status(500).send('Internal Server Error');
    }
});

module.exports = router;
