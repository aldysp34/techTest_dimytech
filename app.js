const express = require('express');
const sqlite3 = require('sqlite3').verbose();

const app = express();
const PORT = 3000;

// Create a new SQLite database connection
const db = new sqlite3.Database(':memory:');

// Middleware to parse JSON requests
app.use(express.json());

// Define API endpoints

// Create a new transaction
app.post('/transactions', (req, res) => {
  const { customer_id, customer_address_id, transaction_date, products, payment_methods } = req.body;

  // Insert transaction into the database
  db.run(
    'INSERT INTO Transaction (customer_id, customer_address_id, transaction_date) VALUES (?, ?, ?)',
    [customer_id, customer_address_id, transaction_date],
    function (err) {
      if (err) {
        console.error(err);
        res.sendStatus(500);
      } else {
        const transactionId = this.lastID;
        // Insert products into the Transaction_Product table
        if (products && Array.isArray(products)) {
          const productValues = products.map((product) => [transactionId, product.product_id, product.quantity]);
          db.run('INSERT INTO Transaction_Product (transaction_id, product_id, quantity) VALUES (?, ?, ?)', productValues, (err) => {
            if (err) {
              console.error(err);
              res.sendStatus(500);
            } else {
              // Insert payment methods into the Transaction_PaymentMethod table
              if (payment_methods && Array.isArray(payment_methods)) {
                const paymentMethodValues = payment_methods.map((paymentMethod) => [transactionId, paymentMethod.payment_method_id]);
                db.run('INSERT INTO Transaction_PaymentMethod (transaction_id, payment_method_id) VALUES (?, ?)', paymentMethodValues, (err) => {
                  if (err) {
                    console.error(err);
                    res.sendStatus(500);
                  } else {
                    res.status(201).json({ transactionId });
                  }
                });
              } else {
                res.status(201).json({ transactionId });
              }
            }
          });
        } else {
          res.status(201).json({ transactionId });
        }
      }
    }
  );
});

// Start the server
app.listen(PORT, () => {
  console.log(`Server listening on port ${PORT}`);
});
