const express = require('express');
const app = express();
const publicRoutes = require('./routes/public.routes');

app.set('view engine', 'ejs');
app.set('views', __dirname + '/views');

app.use('/', publicRoutes);

const PORT = process.env.PORT || 3001;
app.listen(PORT, () => {
    console.log(`Public view service running on port ${PORT}`);
});