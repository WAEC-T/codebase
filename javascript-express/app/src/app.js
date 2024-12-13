const createError = require('http-errors');
const express = require('express');
let path = require('path');
const cookieParser = require('cookie-parser');
const logger = require('morgan');
const flash = require('express-flash');
const session = require('express-session');
const expressLayouts = require('express-ejs-layouts');
const MD5 = require('crypto-js/md5');
const { Users, Messages, Followers, get_user_id, getUserTimelineMessages, validateUser } = require('../../database/sequilize');

const PER_PAGE = 30;

let app = express();

app.set('views', path.join(__dirname, 'views'));
app.set('view engine', 'ejs');

app.use(session({
  secret: process.env.SECRET_KEY,
  resave: false,
  saveUninitialized: true,
}));

app.use(expressLayouts);
app.set('layout', 'layout.ejs');
app.use(logger('dev'));
app.use(flash());
app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(cookieParser());
app.use(express.static(path.join(__dirname, 'public')));

const format_datetime = (timestamp) => {
  return new Date(timestamp * 1000).toISOString().replace(/T/, ' @ ').replace(/\..+/, '');
};

const gravatarUrl = (email, size = 80) => {
  const hash = MD5(email.trim().toLowerCase()).toString();
  return `http://www.gravatar.com/avatar/${hash}?d=identicon&s=${size}`;
};

// Add datetimeformat function to locals object, so it can be called in .ejs views
app.locals.format_datetime = format_datetime;
app.locals.gravatarUrl = gravatarUrl;
app.locals.user = null;

// Routes
app.get('/', async (req, res) => {
  try {
    let user = req.session?.user;
    
    if (!user | !validateUser(user.user_id)) {
      return res.redirect('/public');
    }

    const messages = await getUserTimelineMessages(user.user_id, PER_PAGE);

    console.log("query here ~~~~~~", messages);

    res.render('timeline.ejs', { user: req.session.user, messages, title: 'My Timeline', flashes: req.flash('success'), endpoint: 'my_timeline' });
  } catch (error) {
    console.error('Error: ', error);
    res.status(500).send('Internal Server Error');
  }
});

app.get('/public', async (req, res) => {
  try {
    // Fetch all messages from the database
    const messages = await Messages.findAll({
      order: [['message_id', 'DESC']],
      limit: PER_PAGE,
      raw: true
    });

    // Fetch user details to add to each message
    const user_ids = messages.map(msg => msg.author_id);
    const users = await Users.findAll({
      attributes: ['user_id', 'username', 'email'],
      where: {
        user_id: user_ids
      },
      raw: true
    });

    const usersMap = users.reduce((userMap, user) => {
      userMap[user.user_id] = user;
      return userMap;
    }, {});

    // Add user details to each message
    messages.forEach(msg => {
      msg.username = usersMap[msg.author_id].username;
      msg.email = usersMap[msg.author_id].email;
    });

    // Render the timeline with all messages
    res.render('timeline.ejs', { 
      user: req.session.user, 
      messages: messages, 
      title: 'Public Timeline', 
      flashes: req.flash('success'), 
      endpoint: '' 
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

    if (!username || !email || !password || !password2) {
      return res.render('register.ejs', { error: 'All fields are required', flashes: req.flash('success') });
    }

    if (password !== password2) {
      return res.render('register.ejs', { error: 'The two passwords do not match', flashes: req.flash('success') });
    }

    // Check if email is valid
    if (!email.includes('@')) {
      return res.render('register.ejs', { error: 'You have to enter a valid email address', flashes: req.flash('success') });
    }

    // Check if username is already taken
    const existingUser = await get_user_id(username);
    if (existingUser) {
      return res.render('register.ejs', { error: 'The username is already taken', flashes: req.flash('success') });
    }

    await Users.create({ username: username, email: email, pw_hash: password });

    req.flash('success', 'You were successfully registered and can login now');
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
  res.render('login.ejs', { user: req.session.user, error: null, flashes: req.flash('success') });
});

app.post('/login', async (req, res) => {
  // Implement your logic here
  try {
    const { username, password } = req.body;

    // Query the user from the database
    const user = await Users.findOne({ where: { username: username } });
    if (!user || user.length === 0) {
      return res.render('login.ejs', { error: 'Invalid username', flashes: req.flash('success') });
    }

    if (password !== user.pw_hash) {
      return res.render('login.ejs', { error: 'Invalid password', flashes: req.flash('success') });
    }

    // Set user session
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
  // Check if user is authenticated
  if (!req.session.user) {
    res.status(401).send('Unauthorized');
    return;
  }

  const text = req.body.text;

  if (!text) {
    res.status(400).send('Bad Request');
    return;
  }
  try {
    // Insert Messages into the database
    await Messages.create({
      author_id: req.session.user.user_id,
      text: text,
      pub_date: Math.floor(Date.now() / 1000),
      flagged: 0
    });

    req.flash('success', 'Your message was recorded');
    res.redirect('/');
  } catch (error) {
    console.error('Error inserting message:', error);
    res.status(500).send('Internal Server Error');
  }
});

app.get('/:username', async (req, res) => {
  try {
    // Retrieve the profile user from the databas
    let profile_user = await Users.findOne({
      where: { username: req.params.username }
    });


    // If profile_user is not found, return 404
    if (!profile_user) {
      res.status(404).send('User not found');
      return;
    }
    let followed = false;

    // Check if the logged-in user follows the profile user
    if (req.session.user) {
      var is_following = await Followers.findOne({
        where: { who_id: req.session.user.user_id, whom_id: profile_user.user_id }
      });
      if (is_following) {
        followed = true;
      }
    }

    // Fetch messages for the profile user
    const messages = await Messages.findAll({
      where: {
        author_id: profile_user.user_id
      },
      order: [['pub_date', 'DESC']],
      limit: PER_PAGE,
      raw: true
    });

    // Add user_id, username, and email to each message
    messages.forEach(msg => {
      msg.user_id = profile_user.user_id;
      msg.username = profile_user.username;
      msg.email = profile_user.email;
    });

    // Render the timeline template with messages and other data
    res.render('timeline.ejs', {
      messages: messages,
      followed: followed,
      user: req.session.user, 
      profile_user: profile_user,
      title: profile_user.username + "'s Timeline",
      flashes: req.flash('success'),
      endpoint: 'user_timeline'
    });
  } catch (error) {
    console.error('Error fetching user timeline:', error);
    res.status(500).send('Internal Server Error');
  }
});


app.get('/:username/follow', async (req, res) => {
  try {
    // Check if user is authenticated
    if (!req.session.user) {
      return res.status(401).send('Unauthorized');
    }

    // Get whom_id from the database
    const username = req.params.username;
    const whom_id = await get_user_id(username);
    if (!whom_id) {
      return res.status(404).send('User not found');
    }

    // Insert into follower table
    //await execute('insert into follower (who_id, whom_id) values ($1, $2)', [req.session.user.user_id, whom_id]);
    const newFollower = await Followers.create({
      who_id: req.session.user.user_id,
      whom_id: whom_id
    });

    req.flash('success', `You are now following "${username}"`);
    res.redirect(`/${username}`);
  } catch (error) {
    console.error('Error:', error);
    res.status(500).send('Internal Server Error');
  }
});

app.get('/:username/unfollow', async (req, res) => {
  try {
    // Check if user is authenticated
    if (!req.session.user) {
      return res.status(401).send('Unauthorized');
    }

    // Get whom_id from the database
    const username = req.params.username;
    const whom_id = await get_user_id(username);
    if (!whom_id) {
      return res.status(404).send('User not found');
    }

    // Delete from follower table
    //await execute('delete from follower where who_id=$1 and whom_id=$2', [req.session.user.user_id, whom_id]);
    await Followers.destroy({
      where: {
        who_id: req.session.user.user_id,
        whom_id: whom_id
      }
    });

    req.flash('success', `You are no longer following "${username}"`);
    res.redirect(`/${username}`);
  } catch (error) {
    console.error('Error:', error);
    res.status(500).send('Internal Server Error');
  }
});

// catch 404 and forward to error handler
app.use(function (req, res, next) {
  next(createError(404));
});

app.use(function (err, req, res, next) {
  res.locals.message = err.message || 'An error occurred';
  res.locals.error = req.app.get('env') === 'development' ? err : {};
  res.status(err.status || 500).render('error');
});

app.listen(5000, () => {
  console.log('Minitwit running at port :5000');
});

module.exports = app;