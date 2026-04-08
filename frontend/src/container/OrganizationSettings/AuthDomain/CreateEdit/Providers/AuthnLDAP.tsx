import { Callout } from '@signozhq/callout';
import { Checkbox } from '@signozhq/checkbox';
import { Style } from '@signozhq/design-tokens';
import { CircleHelp } from '@signozhq/icons';
import { Input } from '@signozhq/input';
import { Form, InputNumber, Tooltip } from 'antd';

import RoleMappingSection from './components/RoleMappingSection';

import './Providers.styles.scss';

function ConfigureLDAPAuthnProvider({
	isCreate,
}: {
	isCreate: boolean;
}): JSX.Element {
	const form = Form.useFormInstance();

	return (
		<div className="authn-provider">
			<section className="authn-provider__header">
				<h3 className="authn-provider__title">
					Configure LDAP / Active Directory
				</h3>
				<p className="authn-provider__description">
					Authenticate users against an LDAP or Active Directory server. Users
					will sign in with their email and AD password.
				</p>
			</section>

			<div className="authn-provider__columns">
				{/* Left Column - Core LDAP Settings */}
				<div className="authn-provider__left">
					<div className="authn-provider__field-group">
						<label className="authn-provider__label" htmlFor="ldap-domain">
							Email Domain
							<Tooltip title="The email domain for users who should use LDAP login (e.g., setsoftware.com for users with @setsoftware.com emails)">
								<CircleHelp size={14} color={Style.L3_FOREGROUND} cursor="help" />
							</Tooltip>
						</label>
						<Form.Item
							name="name"
							className="authn-provider__form-item"
							rules={[
								{
									required: true,
									message: 'Email domain is required',
									whitespace: true,
								},
							]}
						>
							<Input id="ldap-domain" disabled={!isCreate} />
						</Form.Item>
					</div>

					<div className="authn-provider__field-group">
						<label className="authn-provider__label" htmlFor="ldap-host">
							LDAP Server Host
							<Tooltip title="The hostname or IP address of your LDAP/AD server (e.g., mail.setyazilim.com.tr)">
								<CircleHelp size={14} color={Style.L3_FOREGROUND} cursor="help" />
							</Tooltip>
						</label>
						<Form.Item
							name={['ldapConfig', 'host']}
							className="authn-provider__form-item"
							rules={[
								{
									required: true,
									message: 'LDAP host is required',
									whitespace: true,
								},
							]}
						>
							<Input id="ldap-host" placeholder="mail.setyazilim.com.tr" />
						</Form.Item>
					</div>

					<div className="authn-provider__field-group">
						<label className="authn-provider__label" htmlFor="ldap-port">
							Port
							<Tooltip title="LDAP server port. Default: 389 for LDAP, 636 for LDAPS. Leave at 0 to use the default.">
								<CircleHelp size={14} color={Style.L3_FOREGROUND} cursor="help" />
							</Tooltip>
						</label>
						<Form.Item
							name={['ldapConfig', 'port']}
							className="authn-provider__form-item"
						>
							<InputNumber
								id="ldap-port"
								min={0}
								max={65535}
								placeholder="389"
								style={{ width: '100%' }}
							/>
						</Form.Item>
					</div>

					<div className="authn-provider__field-group">
						<label className="authn-provider__label" htmlFor="ldap-domains">
							AD Domain Names
							<Tooltip title="Comma-separated Active Directory domain names used for bind. Users authenticate as DOMAIN\username. (e.g., SETYAZILIM,SETSOFTWARE)">
								<CircleHelp size={14} color={Style.L3_FOREGROUND} cursor="help" />
							</Tooltip>
						</label>
						<Form.Item
							name={['ldapConfig', 'domainsInput']}
							className="authn-provider__form-item"
							rules={[
								{
									required: true,
									message: 'At least one AD domain name is required',
									whitespace: true,
								},
							]}
						>
							<Input
								id="ldap-domains"
								placeholder="SETYAZILIM,SETSOFTWARE"
								onChange={(e): void => {
									const domains = e.target.value
										.split(',')
										.map((d: string) => d.trim())
										.filter(Boolean);
									form.setFieldValue(['ldapConfig', 'domains'], domains);
								}}
							/>
						</Form.Item>
					</div>

					<div className="authn-provider__checkbox-row">
						<Form.Item
							name={['ldapConfig', 'useTLS']}
							valuePropName="checked"
							noStyle
						>
							<Checkbox
								id="ldap-use-tls"
								labelName="Use LDAPS (TLS)"
								onCheckedChange={(checked: boolean): void => {
									form.setFieldValue(['ldapConfig', 'useTLS'], checked);
								}}
							/>
						</Form.Item>
						<Tooltip title="Use LDAPS (port 636) instead of plain LDAP (port 389).">
							<CircleHelp size={14} color={Style.L3_FOREGROUND} cursor="help" />
						</Tooltip>
					</div>

					<div className="authn-provider__checkbox-row">
						<Form.Item
							name={['ldapConfig', 'insecureSkipVerify']}
							valuePropName="checked"
							noStyle
						>
							<Checkbox
								id="ldap-skip-verify"
								labelName="Skip TLS Certificate Verification"
								onCheckedChange={(checked: boolean): void => {
									form.setFieldValue(
										['ldapConfig', 'insecureSkipVerify'],
										checked,
									);
								}}
							/>
						</Form.Item>
						<Tooltip title="Skip TLS certificate verification. Enable this for self-signed certificates.">
							<CircleHelp size={14} color={Style.L3_FOREGROUND} cursor="help" />
						</Tooltip>
					</div>

					<Callout
						type="info"
						size="small"
						showIcon
						description="Users will authenticate as DOMAIN\username against your AD server. The email prefix is used as the username."
						className="callout"
					/>
				</div>

				{/* Right Column - Optional Settings */}
				<div className="authn-provider__right">
					<div className="authn-provider__field-group">
						<label className="authn-provider__label" htmlFor="ldap-base-dn">
							Base DN (optional)
							<Tooltip title="Base DN for user search. If empty, only bind authentication is performed. (e.g., DC=setyazilim,DC=com,DC=tr)">
								<CircleHelp size={14} color={Style.L3_FOREGROUND} cursor="help" />
							</Tooltip>
						</label>
						<Form.Item
							name={['ldapConfig', 'baseDN']}
							className="authn-provider__form-item"
						>
							<Input
								id="ldap-base-dn"
								placeholder="DC=setyazilim,DC=com,DC=tr"
							/>
						</Form.Item>
					</div>

					<div className="authn-provider__field-group">
						<label
							className="authn-provider__label"
							htmlFor="ldap-user-search-filter"
						>
							User Search Filter (optional)
							<Tooltip title="LDAP filter to search for users. Use %s as placeholder for username. Default: (sAMAccountName=%s)">
								<CircleHelp size={14} color={Style.L3_FOREGROUND} cursor="help" />
							</Tooltip>
						</label>
						<Form.Item
							name={['ldapConfig', 'userSearchFilter']}
							className="authn-provider__form-item"
						>
							<Input
								id="ldap-user-search-filter"
								placeholder="(sAMAccountName=%s)"
							/>
						</Form.Item>
					</div>

					<RoleMappingSection
						fieldNamePrefix={['roleMapping']}
						isExpanded={false}
						onExpandChange={(): void => {}}
					/>
				</div>
			</div>
		</div>
	);
}

export default ConfigureLDAPAuthnProvider;
