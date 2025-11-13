import React from "react";
import Layout from "../components/layout";
import { Link } from "gatsby";

const PrivacyPolicyPage: React.FC = () => {
  return (
    <Layout>
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">
            Privacy Policy
          </h1>
          <p className="text-sm text-gray-500">
            Last Updated: {new Date().toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })}
          </p>
        </div>

        <div className="prose prose-lg max-w-none">
          <section className="mb-8">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">
              1. Introduction
            </h2>
            <p className="text-gray-700 mb-4">
              Welcome to EgalDeutsch. We respect your privacy and are committed
              to protecting your personal data. This Privacy Policy explains how
              we collect, use, disclose, and safeguard your information when you
              use our German learning platform.
            </p>
            <p className="text-gray-700 mb-4">
              This Privacy Policy complies with the General Data Protection
              Regulation (GDPR) of the European Union and German data protection
              laws (BDSG - Bundesdatenschutzgesetz).
            </p>
            <div className="bg-blue-50 rounded-lg p-6">
              <p className="text-gray-700">
                <strong>Data Controller:</strong> Steve Phan
              </p>
              <p className="text-gray-700">
                <strong>Contact Email:</strong>{" "}
                <a
                  href="mailto:egaldeutsch.com@gmail.com"
                  className="text-blue-600 hover:underline"
                >
                  egaldeutsch.com@gmail.com
                </a>
              </p>
            </div>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">
              2. Information We Collect
            </h2>
            
            <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-6">
              2.1 Information You Provide
            </h3>
            <p className="text-gray-700 mb-4">
              When you register for an account or use our Service, we may
              collect:
            </p>
            <ul className="list-disc pl-6 text-gray-700 space-y-2">
              <li>Name (first name and last name)</li>
              <li>Email address</li>
              <li>Username</li>
              <li>Password (encrypted)</li>
              <li>Language proficiency level</li>
              <li>Learning preferences and progress</li>
            </ul>

            <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-6">
              2.2 Automatically Collected Information
            </h3>
            <p className="text-gray-700 mb-4">
              When you access our Service, we automatically collect:
            </p>
            <ul className="list-disc pl-6 text-gray-700 space-y-2">
              <li>IP address</li>
              <li>Browser type and version</li>
              <li>Device information</li>
              <li>Operating system</li>
              <li>Referring/exit pages</li>
              <li>Date and time stamps</li>
              <li>Clickstream data</li>
            </ul>

            <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-6">
              2.3 Learning and Usage Data
            </h3>
            <ul className="list-disc pl-6 text-gray-700 space-y-2">
              <li>Quiz scores and performance</li>
              <li>Stories read and completion rates</li>
              <li>Time spent on the platform</li>
              <li>Learning progress and statistics</li>
              <li>Leaderboard rankings</li>
            </ul>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">
              3. How We Use Your Information
            </h2>
            <p className="text-gray-700 mb-4">
              We use the collected information for the following purposes:
            </p>
            <ul className="list-disc pl-6 text-gray-700 space-y-2">
              <li>
                <strong>Service Provision:</strong> To provide and maintain our
                learning platform
              </li>
              <li>
                <strong>Account Management:</strong> To create and manage your
                user account
              </li>
              <li>
                <strong>Personalization:</strong> To personalize your learning
                experience
              </li>
              <li>
                <strong>Progress Tracking:</strong> To track your learning
                progress and achievements
              </li>
              <li>
                <strong>Communication:</strong> To send important updates,
                notifications, and educational content
              </li>
              <li>
                <strong>Analytics:</strong> To analyze usage patterns and improve
                our Service
              </li>
              <li>
                <strong>Security:</strong> To detect, prevent, and address
                technical issues and security threats
              </li>
              <li>
                <strong>Legal Compliance:</strong> To comply with legal
                obligations
              </li>
            </ul>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">
              4. Legal Basis for Processing (GDPR)
            </h2>
            <p className="text-gray-700 mb-4">
              Under GDPR, we process your personal data based on:
            </p>
            <ul className="list-disc pl-6 text-gray-700 space-y-2">
              <li>
                <strong>Consent:</strong> You have given consent for specific
                purposes (Article 6(1)(a) GDPR)
              </li>
              <li>
                <strong>Contract Performance:</strong> Processing is necessary
                for providing our Service (Article 6(1)(b) GDPR)
              </li>
              <li>
                <strong>Legal Obligation:</strong> Processing is required by law
                (Article 6(1)(c) GDPR)
              </li>
              <li>
                <strong>Legitimate Interests:</strong> Processing is necessary
                for our legitimate interests, such as improving our Service
                (Article 6(1)(f) GDPR)
              </li>
            </ul>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">
              5. Third-Party Services
            </h2>
            <p className="text-gray-700 mb-4">
              We use the following third-party services that may process your
              data:
            </p>

            <div className="space-y-4">
              <div className="bg-gray-50 rounded-lg p-6">
                <h3 className="text-lg font-semibold text-gray-900 mb-2">
                  Netlify (Hosting)
                </h3>
                <p className="text-gray-700 mb-2">
                  Our platform is hosted on Netlify's infrastructure. Netlify may
                  process technical data such as IP addresses and usage logs.
                </p>
                <p className="text-gray-700 text-sm">
                  Privacy Policy:{" "}
                  <a
                    href="https://www.netlify.com/privacy/"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-blue-600 hover:underline"
                  >
                    https://www.netlify.com/privacy/
                  </a>
                </p>
              </div>

              <div className="bg-gray-50 rounded-lg p-6">
                <h3 className="text-lg font-semibold text-gray-900 mb-2">
                  Google Services (Analytics & Functionality)
                </h3>
                <p className="text-gray-700 mb-2">
                  We may use Google Analytics or other Google services to analyze
                  how users interact with our platform. Google may use cookies
                  and similar technologies to collect usage data.
                </p>
                <p className="text-gray-700 text-sm">
                  Privacy Policy:{" "}
                  <a
                    href="https://policies.google.com/privacy"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-blue-600 hover:underline"
                  >
                    https://policies.google.com/privacy
                  </a>
                </p>
              </div>

              <div className="bg-gray-50 rounded-lg p-6">
                <h3 className="text-lg font-semibold text-gray-900 mb-2">
                  MongoDB (Database)
                </h3>
                <p className="text-gray-700 mb-2">
                  User data and learning progress are securely stored in MongoDB
                  databases with encryption and access controls.
                </p>
                <p className="text-gray-700 text-sm">
                  Privacy Policy:{" "}
                  <a
                    href="https://www.mongodb.com/legal/privacy-policy"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-blue-600 hover:underline"
                  >
                    https://www.mongodb.com/legal/privacy-policy
                  </a>
                </p>
              </div>
            </div>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">
              6. Cookies and Tracking Technologies
            </h2>
            <p className="text-gray-700 mb-4">
              We use cookies and similar tracking technologies to enhance your
              experience. Cookies are small data files stored on your device.
            </p>
            <p className="text-gray-700 mb-4">Types of cookies we use:</p>
            <ul className="list-disc pl-6 text-gray-700 space-y-2">
              <li>
                <strong>Essential Cookies:</strong> Required for the Service to
                function (e.g., authentication)
              </li>
              <li>
                <strong>Analytics Cookies:</strong> Help us understand how users
                interact with our Service
              </li>
              <li>
                <strong>Preference Cookies:</strong> Remember your settings and
                preferences
              </li>
            </ul>
            <p className="text-gray-700 mt-4">
              You can control cookie settings through your browser, but disabling
              certain cookies may limit functionality.
            </p>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">
              7. Data Retention
            </h2>
            <p className="text-gray-700 mb-4">
              We retain your personal data only as long as necessary for the
              purposes outlined in this Privacy Policy:
            </p>
            <ul className="list-disc pl-6 text-gray-700 space-y-2">
              <li>
                <strong>Account Data:</strong> Retained while your account is
                active and for 30 days after account deletion
              </li>
              <li>
                <strong>Learning Data:</strong> Retained to track your progress
                unless you request deletion
              </li>
              <li>
                <strong>Analytics Data:</strong> Anonymized and retained for
                statistical purposes
              </li>
            </ul>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">
              8. Your Rights Under GDPR
            </h2>
            <p className="text-gray-700 mb-4">
              As a user in the EU/Germany, you have the following rights:
            </p>
            <ul className="list-disc pl-6 text-gray-700 space-y-3">
              <li>
                <strong>Right of Access:</strong> Request a copy of your personal
                data we hold
              </li>
              <li>
                <strong>Right to Rectification:</strong> Correct inaccurate or
                incomplete personal data
              </li>
              <li>
                <strong>Right to Erasure ("Right to be Forgotten"):</strong>{" "}
                Request deletion of your personal data
              </li>
              <li>
                <strong>Right to Restriction:</strong> Limit how we use your
                personal data
              </li>
              <li>
                <strong>Right to Data Portability:</strong> Receive your data in
                a structured, machine-readable format
              </li>
              <li>
                <strong>Right to Object:</strong> Object to processing of your
                personal data
              </li>
              <li>
                <strong>Right to Withdraw Consent:</strong> Withdraw your consent
                at any time
              </li>
              <li>
                <strong>Right to Lodge a Complaint:</strong> File a complaint
                with a supervisory authority
              </li>
            </ul>
            <p className="text-gray-700 mt-4">
              To exercise any of these rights, please contact us at{" "}
              <a
                href="mailto:egaldeutsch.com@gmail.com"
                className="text-blue-600 hover:underline"
              >
                egaldeutsch.com@gmail.com
              </a>
            </p>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">
              9. Data Security
            </h2>
            <p className="text-gray-700 mb-4">
              We implement appropriate technical and organizational measures to
              protect your personal data:
            </p>
            <ul className="list-disc pl-6 text-gray-700 space-y-2">
              <li>Encryption of data in transit (HTTPS/TLS)</li>
              <li>Encryption of sensitive data at rest</li>
              <li>Password hashing using industry-standard algorithms</li>
              <li>Regular security assessments and updates</li>
              <li>Access controls and authentication mechanisms</li>
              <li>Secure database configurations</li>
            </ul>
            <p className="text-gray-700 mt-4">
              However, no method of transmission over the Internet is 100%
              secure. While we strive to protect your data, we cannot guarantee
              absolute security.
            </p>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">
              10. International Data Transfers
            </h2>
            <p className="text-gray-700">
              Your data may be transferred to and processed in countries outside
              the EU/EEA. When we transfer data internationally, we ensure
              appropriate safeguards are in place, such as Standard Contractual
              Clauses (SCCs) approved by the European Commission, to protect your
              data in accordance with GDPR requirements.
            </p>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">
              11. Children's Privacy
            </h2>
            <p className="text-gray-700">
              Our Service is intended for users aged 16 and above. If you are
              under 16, please obtain parental consent before using our Service.
              We do not knowingly collect personal data from children under 16
              without parental consent. If we become aware that we have collected
              personal data from a child under 16 without parental consent, we
              will take steps to delete such information.
            </p>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">
              12. Changes to This Privacy Policy
            </h2>
            <p className="text-gray-700">
              We may update this Privacy Policy from time to time. We will notify
              you of any material changes by posting the new Privacy Policy on
              this page and updating the "Last Updated" date. We encourage you to
              review this Privacy Policy periodically.
            </p>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">
              13. Supervisory Authority
            </h2>
            <p className="text-gray-700 mb-4">
              If you have concerns about how we handle your personal data, you
              have the right to lodge a complaint with a data protection
              supervisory authority.
            </p>
            <div className="bg-gray-50 rounded-lg p-6">
              <p className="text-gray-700 mb-2">
                <strong>For Germany:</strong>
              </p>
              <p className="text-gray-700">
                Die Bundesbeauftragte für den Datenschutz und die
                Informationsfreiheit (BfDI)
                <br />
                <a
                  href="https://www.bfdi.bund.de"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-blue-600 hover:underline"
                >
                  https://www.bfdi.bund.de
                </a>
              </p>
            </div>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">
              14. Contact Us
            </h2>
            <p className="text-gray-700 mb-4">
              If you have any questions about this Privacy Policy or our data
              practices, please contact us:
            </p>
            <div className="bg-gray-50 rounded-lg p-6">
              <p className="text-gray-700">
                <strong>Data Controller:</strong> Steve Phan
              </p>
              <p className="text-gray-700">
                <strong>Email:</strong>{" "}
                <a
                  href="mailto:egaldeutsch.com@gmail.com"
                  className="text-blue-600 hover:underline"
                >
                  egaldeutsch.com@gmail.com
                </a>
              </p>
            </div>
          </section>
        </div>

        <div className="mt-12 pt-8 border-t border-gray-200">
          <div className="text-center space-x-6">
            <Link
              to="/terms-of-service"
              className="text-blue-600 hover:underline font-medium"
            >
              Terms of Service
            </Link>
            <Link
              to="/"
              className="text-blue-600 hover:underline font-medium"
            >
              Back to Home
            </Link>
          </div>
        </div>
      </div>
    </Layout>
  );
};

export default PrivacyPolicyPage;
