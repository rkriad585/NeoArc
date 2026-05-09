def register_route_blueprints(app):
    from core.routes.auth_routes import auth_bp
    from core.routes.alias_routes import alias_bp
    from core.routes.profile_routes import profile_bp
    from core.routes.api_routes import api_bp
    app.register_blueprint(auth_bp)
    app.register_blueprint(alias_bp)
    app.register_blueprint(profile_bp)
    app.register_blueprint(api_bp)
