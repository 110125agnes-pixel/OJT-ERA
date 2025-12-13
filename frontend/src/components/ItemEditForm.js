import React from 'react';

const ItemEditForm = ({ item, onChange, onSave, onCancel, loading, sexOptions, civilStatusOptions }) => {
  return (
    <>
      <div className="form-row">
        <input
          type="text"
          value={item.lastname}
          onChange={e => onChange({ ...item, lastname: e.target.value })}
          placeholder="Lastname"
          disabled={loading}
          required
        />
        <input
          type="text"
          value={item.firstname}
          onChange={e => onChange({ ...item, firstname: e.target.value })}
          placeholder="Firstname"
          disabled={loading}
          required
        />
        <input
          type="text"
          value={item.middlename}
          onChange={e => onChange({ ...item, middlename: e.target.value })}
          placeholder="Middlename"
          disabled={loading}
        />
        <input
          type="text"
          value={item.suffix}
          onChange={e => onChange({ ...item, suffix: e.target.value })}
          placeholder="Suffix"
          disabled={loading}
        />
      </div>
      <div className="form-row">
        <input
          type="date"
          value={item.birthdate}
          onChange={e => onChange({ ...item, birthdate: e.target.value })}
          placeholder="Birthdate"
          disabled={loading}
        />
        <select
          value={item.sex}
          onChange={e => onChange({ ...item, sex: e.target.value })}
          disabled={loading}
          required
        >
          {sexOptions.map(opt => <option key={opt} value={opt}>{opt || 'Sex'}</option>)}
        </select>
        <select
          value={item.civil_status}
          onChange={e => onChange({ ...item, civil_status: e.target.value })}
          disabled={loading}
          required
        >
          {civilStatusOptions.map(opt => <option key={opt} value={opt}>{opt || 'Civil Status'}</option>)}
        </select>
      </div>
      <div className="button-group">
        <button onClick={onSave} disabled={loading} className="save-btn">Save</button>
        <button onClick={onCancel} disabled={loading} className="cancel-btn">Cancel</button>
      </div>
    </>
  );
};

export default ItemEditForm;
