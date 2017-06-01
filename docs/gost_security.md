## GOST and security

The SensorThings API does not define specific security capabilities, security aspects should be handled
in a different application layer (see http://docs.opengeospatial.org/is/15-078r6/15-078r6.html#21).

One common usecase is to limit write access to the SensorThings API. By configuring a proxy server like Apache2 or 
NGINX this can be achieved.

Example configuration for Apache2:

```
<Location /gost/v1.0/>
		<Limit POST PUT DELETE PATCH>
			AuthType Basic
			AuthName "GOST Writer"
			AuthUserFile "/etc/apache2/thepwfile"
			Require user gost
		</Limit>
	</Location>
```

And for Dashboard:

```
<Location "/gost/Dashboard/">
		AuthType Basic
		AuthName "GOST Writer"
		AuthUserFile "/etc/apache2/thepwfile"
		Require user gost
	</Location>
  ```
  
  Example configuration for NGINX:
  
  TODO


