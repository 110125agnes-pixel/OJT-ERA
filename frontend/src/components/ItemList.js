import React from 'react';

const ItemList = ({ items, onEdit, onDelete, loading }) => {
  if (loading && items.length === 0) {
    return <p>Loading...</p>;
  }

  if (items.length === 0) {
    return <p>No items yet. Add one above!</p>;
  }

  return (
    <ul>
      {items.map((item) => (
        <li key={item.id}>
          <div className="item-fields">
            <span><b>{item.lastname}, {item.firstname} {item.middlename} {item.suffix}</b></span><br />
            <span>Birthdate: {item.birthdate || '-'}</span> | <span>Sex: {item.sex || '-'}</span> | <span>Civil Status: {item.civil_status || '-'}</span>
          </div>
          <div className="button-group">
            <button onClick={() => onEdit(item)} disabled={loading} className="edit-btn">Edit</button>
            <button onClick={() => onDelete(item.id)} disabled={loading} className="delete-btn">Delete</button>
          </div>
        </li>
      ))}
    </ul>
  );
};

export default ItemList;
