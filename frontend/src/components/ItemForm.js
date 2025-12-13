import React from 'react';

const ItemForm = ({ formData, onChange, onSubmit, loading, sexOptions, civilStatusOptions }) => {
  return (
    <form onSubmit={onSubmit} className="add-form">
      <div className="form-row">
        <input
          type="text"
          value={formData.lastname}
          onChange={e => onChange({ ...formData, lastname: e.target.value })}
          placeholder="Lastname"
          disabled={loading}
          required
        />
        <input
          type="text"
          value={formData.firstname}
          onChange={e => onChange({ ...formData, firstname: e.target.value })}
          placeholder="Firstname"
          disabled={loading}
          required
        />
        <input
          type="text"
          value={formData.middlename}
          onChange={e => onChange({ ...formData, middlename: e.target.value })}
          placeholder="Middlename"
          disabled={loading}
        />
        <input
          type="text"
          value={formData.suffix}
          onChange={e => onChange({ ...formData, suffix: e.target.value })}
          placeholder="Suffix"
          disabled={loading}
        />
      </div>
      <div className="form-row">
        <input
          type="date"
          value={formData.birthdate}
          onChange={e => onChange({ ...formData, birthdate: e.target.value })}
          placeholder="Birthdate"
          disabled={loading}
        />
        <select
          value={formData.sex}
          onChange={e => onChange({ ...formData, sex: e.target.value })}
          disabled={loading}
          required
        >
          {sexOptions.map(opt => <option key={opt} value={opt}>{opt || 'Sex'}</option>)}
        </select>
        <select
          value={formData.civil_status}
          onChange={e => onChange({ ...formData, civil_status: e.target.value })}
          disabled={loading}
          required
        >
          {civilStatusOptions.map(opt => <option key={opt} value={opt}>{opt || 'Civil Status'}</option>)}
        </select>
      </div>
      <button type="submit" disabled={loading}>
        {loading ? 'Adding...' : 'Add Item'}
      </button>
    </form>
  );
};

export default ItemForm;
