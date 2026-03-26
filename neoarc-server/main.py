# main.py

from flask import Flask, request, jsonify, render_template, redirect, url_for, session, flash
from werkzeug.security import generate_password_hash, check_password_hash

from core.admin import admin_bp
from core.db import init_db, get_db

import config
import sqlite3
import os

app = Flask(__name__)
app.register_blueprint(admin_bp, url_prefix=config.ADMIN_ROUTE)
app.secret_key = config.SECRET_KEY

app.config['UPLOAD_FOLDER'] = config.UPLOAD_FOLDER

# Ensure upload directory exists
if not os.path.exists(config.UPLOAD_FOLDER):
    os.makedirs(config.UPLOAD_FOLDER)

def allowed_file(filename):
    return '.' in filename and filename.rsplit('.', 1)[1].lower() in config.ALLOWED_EXTENSIONS

init_db()

if not os.path.exists(config.DB_PATH):
    init_db()

@app.errorhandler(404)
def page_not_found(e):
    return render_template('error.html', code="404", message="NODE_NOT_FOUND", subtext="The requested sector has been purged or does not exist in the matrix."), 404

@app.errorhandler(500)
def internal_error(e):
    return render_template('error.html', code="500", message="SYSTEM_FAILURE", subtext="A critical error occurred within the mainframe. Engineers have been notified."), 500

@app.route('/')
def index():
    if 'user_id' in session:
        return redirect(url_for('dashboard'))
    return render_template('login.html')

@app.route('/login', methods=['POST'])
def login():
    username, password = request.form['username'], request.form['password']
    conn = get_db()
    user = conn.execute("SELECT * FROM users WHERE username=?", (username,)).fetchone()
    if user and check_password_hash(user['password'], password):
        
        if user['is_blocked'] == 1:
            flash("ACCESS DENIED: Your account has been blocked by an Administrator.", "error")
            return redirect(url_for('index'))
        
        session['user_id'] = user['id']
        session['username'] = user['username']
        session['profile_pic'] = user['profile_pic'] if user['profile_pic'] else None
        flash("Authentication successful. Welcome to the NeoArc.", "success")
        return redirect(url_for('dashboard'))
    flash("Invalid credentials, matrix denied.", "error")
    return redirect(url_for('index'))

@app.route('/register', methods=['GET', 'POST'])
def register():
    if request.method == 'POST':
        full_name = request.form['full_name']
        email = request.form['email']
        username = request.form['username']
        password = generate_password_hash(request.form['password'])
        
        try:
            conn = get_db()
            conn.execute("INSERT INTO users (username, email, full_name, password) VALUES (?, ?, ?, ?)", 
                         (username, email, full_name, password))
            conn.commit()
            flash("Identity established. You may now login.", "success")
            return redirect(url_for('index'))
        except sqlite3.IntegrityError:
            flash("Username or Email already exists in the grid.", "error")
            
    return render_template('register.html')

@app.route('/dashboard', methods=['GET', 'POST'])
def dashboard():
    if 'user_id' not in session: return redirect(url_for('index'))
    conn = get_db()
    
    if request.method == 'POST': # Create Alias
        alias_name = request.form['alias_name']
        command = request.form['command']
        exec_type = request.form.get('exec_type', '') 
        try:
            conn.execute("INSERT INTO aliases (alias_name, command, exec_type, user_id) VALUES (?, ?, ?, ?)", 
                         (alias_name, command, exec_type, session['user_id']))
            conn.commit()
            flash("Alias successfully Created into NeoArc.", "success")
        except sqlite3.IntegrityError:
            flash("Alias name already exists!", "error") # Handle duplicate alias error here
            
    aliases = conn.execute("SELECT * FROM aliases WHERE user_id=?", (session['user_id'],)).fetchall()
    return render_template('dashboard.html', aliases=aliases, username=session['username'])

@app.route('/profile', methods=['GET', 'POST'])
def profile():
    if 'user_id' not in session: return redirect(url_for('index'))
    conn = get_db()
    
    if request.method == 'POST':
        full_name = request.form['full_name']
        email = request.form['email']

        # Handle Image Upload
        file = request.files.get('profile_pic')
        filename = None
        if file and file.filename != '' and allowed_file(file.filename):
            ext = file.filename.rsplit('.', 1)[1].lower()
            # Rename file to user_id.ext to prevent duplicates/overwrites
            filename = f"user_{session['user_id']}.{ext}"
            file.save(os.path.join(app.config['UPLOAD_FOLDER'], filename))
            
            # Update Session immediately
            session['profile_pic'] = filename
            
            # Update DB with image
            conn.execute("UPDATE users SET full_name=?, email=?, profile_pic=? WHERE id=?", 
                         (full_name, email, filename, session['user_id']))
        else:
            # Update DB without image
            conn.execute("UPDATE users SET full_name=?, email=? WHERE id=?", 
                         (full_name, email, session['user_id']))
            
        conn.commit()
        flash("Profile updated successfully.", "success")

        try:
            conn.execute("UPDATE users SET full_name=?, email=? WHERE id=?", (full_name, email, session['user_id']))
            conn.commit()
            flash("Identity profile updated successfully.", "success")
        except sqlite3.IntegrityError:
            flash("Email is already bound to another node.", "error")
            
    user = conn.execute("SELECT * FROM users WHERE id=?", (session['user_id'],)).fetchone()
    node_count = conn.execute("SELECT COUNT(*) FROM aliases WHERE user_id=?", (session['user_id'],)).fetchone()[0]
    
    return render_template('profile.html', user=user, node_count=node_count)

@app.route('/search')
def search():
    if 'user_id' not in session: return redirect(url_for('index'))
    query = request.args.get('q', '')
    
    conn = get_db()
    if query:
        # Search global aliases and join with user info
        results = conn.execute('''
            SELECT a.*, u.username, u.profile_pic 
            FROM aliases a 
            JOIN users u ON a.user_id = u.id 
            WHERE a.alias_name LIKE ? OR a.command LIKE ?
            ORDER BY a.id DESC
        ''', ('%' + query + '%', '%' + query + '%')).fetchall()
    else:
        results = []
        
    return render_template('search.html', results=results, query=query)

@app.route('/view/<alias_name>')
def view_alias(alias_name):
    if 'user_id' not in session: return redirect(url_for('index'))
    conn = get_db()
    
    # Fetch alias details + Creator info
    alias = conn.execute('''
        SELECT a.*, u.username, u.profile_pic, u.email
        FROM aliases a 
        JOIN users u ON a.user_id = u.id 
        WHERE a.alias_name = ?
    ''', (alias_name,)).fetchone()
    
    if not alias:
        return render_template('error.html', code="404", message="ALIAS_MISSING", subtext="This node does not exist in the public matrix."), 404
        
    return render_template('view.html', alias=alias)

@app.route('/delete/<int:id>')
def delete_alias(id):
    if 'user_id' not in session: return redirect(url_for('index'))
    conn = get_db()
    conn.execute("DELETE FROM aliases WHERE id=? AND user_id=?", (id, session['user_id']))
    conn.commit()
    flash("Alias purged from the system.", "success")
    return redirect(url_for('dashboard'))

@app.route('/logout')
def logout():
    session.clear()
    flash("Session terminated.", "success")
    return redirect(url_for('index'))

@app.route('/api/alias/<alias_name>', methods=['GET'])
def get_alias(alias_name):
    result = get_db().execute("SELECT command, exec_type FROM aliases WHERE alias_name=?", (alias_name,)).fetchone()
    if result: 
        return jsonify({"success": True, "alias": alias_name, "command": result['command'], "exec_type": result['exec_type']}), 200
    return jsonify({"success": False, "message": "Alias not found"}), 404

if __name__ == '__main__':
    app.run(host=config.HOST , port=config.PORT, debug=config.DEBUG)