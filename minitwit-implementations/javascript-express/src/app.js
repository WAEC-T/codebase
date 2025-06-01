const express = require('express');
let path = require('path');
const cookieParser = require('cookie-parser');
const flash = require('express-flash');
const session = require('express-session');
const expressLayouts = require('express-ejs-layouts');
const MD5 = require('crypto-js/md5');
const page = require('./routes/page');
const api = require('./routes/api');

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

// Add functions as locals object to be called in .ejs views
app.locals.formatDatetime = formatDatetime;
app.locals.gravatarUrl = gravatarUrl;
app.locals.user = null;

app.use('/', page);
app.use('/api', api);

app.listen(5000, () => {
    console.log('Minitwit running at port :5000');
});

module.exports = app;
