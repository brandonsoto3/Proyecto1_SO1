'use strict';

require('dotenv').config();

var logger = require('morgan');
const express = require('express');
const port = process.env.PORT || 4000;
const path = require('path');

const app = express();

app.use(logger('dev'))
app.engine('html', require('ejs').renderFile);
app.set('views', path.join(__dirname));

app.use(express.static(path.join(__dirname)));


app.use('/', function(req, res, next) {
    res.render('index.html');
});


app.listen(port, function() {
    console.log('Servidor corriendo en el puerto ' + port);
});