const createError = require('http-errors');
const express = require('express');
let path = require('path');
const cookieParser = require('cookie-parser');
const logger = require('morgan');
const flash = require('express-flash');
const session = require('express-session');
const expressLayouts = require('express-ejs-layouts');
const MD5 = require('crypto-js/md5');
const {
    Users,
    Messages,
    Followers,
    getUserIdByName,
    getUserTimelineMessages,
    getPublicMessages,
    getUser,
    getUserByName,
    isFollowing,
    getUserMessages,
    followUser,
    unfollowUser,
} = require('../../database/sequilize');
const { formatMessages, validateRegisterFields } = require('./utils');

const PER_PAGE = 30;

let app = express();

app.set('views', path.join(__dirname, 'views'));
app.set('view engine', 'ejs');

app.use(
    session({
        secret: process.env.SECRET_KEY,
        resave: false,
        saveUninitialized: true,
    })
);

app.use(expressLayouts);
app.set('layout', 'layout.ejs');
app.use(logger('dev'));
app.use(flash());
app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(cookieParser());
app.use(express.static(path.join(__dirname, 'public')));

const formatDatetime = (date) => {
    return new Date(date).toISOString().replace(/T/, ' @ ').replace(/\..+/, '');
};

const gravatarUrl = (email, size = 80) => {
    const hash = MD5(email.trim().toLowerCase()).toString();
    return `http://www.gravatar.com/avatar/${hash}?d=identicon&s=${size}`;
};

// Add datetimeformat function to locals object, so it can be called in .ejs views
app.locals.formatDatetime = formatDatetime;
app.locals.gravatarUrl = gravatarUrl;
app.locals.user = null;

// Routes
app.get('/', async (req, res) => {
    try {
        let user = req.session?.user;

        if (!user || !getUser(user.user_id)) {
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
            endpoint: 'my_timeline',
        });
    } catch (error) {
        console.error('Error: ', error);

        res.status(500).send('Internal Server Error');
    }
});

app.get('/public', async (req, res) => {
    try {
        let user = req.session.user;

        const userRefresh = user ? await getUser(user.user_id) : user;

        const messages = formatMessages(await getPublicMessages(PER_PAGE));

        res.render('timeline.ejs', {
            user: userRefresh,
            messages: messages,
            title: 'Public Timeline',
            flashes: req.flash('success'),
            endpoint: 'public_timeline',
        });
    } catch (error) {
        console.error('Error:', error);

        res.status(500).send('Internal Server Error');
    }
});

app.get('/register', async (req, res) => {
    if (req.session.user) {
        return res.redirect('/timeline');
    }

    res.render('register.ejs', { error: null, flashes: req.flash('success') });
});

app.post('/register', async (req, res) => {
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

        await Users.create({
            username: username,
            email: email,
            pw_hash: password,
        });

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

app.get('/login', (req, res) => {
    if (req.session.user) {
        return res.redirect('/timeline');
    }

    res.render('login.ejs', {
        user: req.session.user,
        error: null,
        flashes: req.flash('success'),
    });
});

app.post('/login', async (req, res) => {
    try {
        const { username, password } = req.body;

        if (!username || !password) {
            return res.render('login.ejs', {
                error: 'Empty fields not allowed',
                flashes: req.flash('success'),
            });
        }

        const user = await Users.findOne({ where: { username: username } });

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

app.get('/logout', async (req, res) => {
    req.flash('success', 'You were logged out');
    delete req.session.user;
    res.redirect('/public');
});

app.post('/add_message', async (req, res) => {
    if (!req.session.user) {
        res.status(401).send('Unauthorized');
        return;
    }

    const text = req.body.text;

    if (!text) {
        req.flash('success', 'Message cannot be empty!'); // minimal effort maximum result ¯\_(ツ)_/¯

        res.redirect('/');

        return;
    }

    try {
        await Messages.create({
            author_id: req.session.user.user_id,
            text: text,
            pub_date: new Date(),
            flagged: 0,
        });

        req.flash('success', 'Your message was recorded');

        res.redirect('/');
    } catch (error) {
        console.error('Error inserting message:', error);

        res.status(500).send('Internal Server Error');
    }
});

app.get('/user/:username', async (req, res) => {
    try {
        let user = req.session?.user;

        const userRefresh = user ? await getUser(user.user_id) : user;

        let profileUser = await getUserByName(req.params.username);

        if (!profileUser) {
            res.status(404).send('User not found');

            return;
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

app.get('/:username/follow', async (req, res) => {
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
            req.flash('success', `You are now following "${whomUsername}"`);

        res.redirect(`/user/${whomUsername}`);

    } catch (error) {
        console.error('Error:', error);

        res.status(500).send('Internal Server Error');
    }
});

app.get('/:username/unfollow', async (req, res) => {
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
            req.flash(
                'success',
                `You are no longer following "${whomUsername}"`
            );

        res.redirect(`/user/${whomUsername}`);

    } catch (error) {
        console.error('Error:', error);

        res.status(500).send('Internal Server Error');
    }
});

app.listen(5000, () => {
    console.log('Minitwit running at port :5000');
});

module.exports = app;
