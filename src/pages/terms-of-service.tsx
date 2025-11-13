import React from "react";
import Layout from "../components/layout";
import { Link } from "gatsby";

const TermsOfServicePage: React.FC = () => {
  return (
    <Layout>
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">
            Terms of Service
          </h1>
          <p className="text-sm text-gray-500">
            Last Updated: {new Date().toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })}
          </p>
        </div>

        <div className="prose prose-lg max-w-none">
          <section className="mb-8">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">
              1. Acceptance of Terms
            </h2>
            <p className="text-gray-700 mb-4">
              Welcome to EgalDeutsch. By accessing or using our German learning
              platform ("Service"), you agree to be bound by these Terms of
              Service ("Terms"). If you do not agree to these Terms, please do
              not use our Service.
            </p>
            <p className="text-gray-700">
              These Terms constitute a legally binding agreement between you and
              EgalDeutsch, operated by Steve Phan. By using our Service, you
              represent that you are at least 16 years old or have parental
              consent if you are younger.
            </p>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">
              2. Description of Service
            </h2>
            <p className="text-gray-700 mb-4">
              EgalDeutsch provides an online platform for learning German
              language through:
            </p>
            <ul className="list-disc pl-6 text-gray-700 space-y-2">
              <li>Reading German stories at various proficiency levels (A1-C2)</li>
              <li>Interactive quizzes to test comprehension</li>
              <li>Vocabulary learning tools</li>
              <li>Progress tracking and leaderboards</li>
              <li>User accounts and personalized learning experiences</li>
            </ul>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">
              3. User Accounts
            </h2>
            <p className="text-gray-700 mb-4">
              To access certain features of our Service, you may be required to
              create an account. When creating an account, you agree to:
            </p>
            <ul className="list-disc pl-6 text-gray-700 space-y-2">
              <li>Provide accurate, current, and complete information</li>
              <li>Maintain the security of your password and account</li>
              <li>Notify us immediately of any unauthorized use of your account</li>
              <li>Accept responsibility for all activities that occur under your account</li>
            </ul>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">
              4. Acceptable Use
            </h2>
            <p className="text-gray-700 mb-4">You agree NOT to:</p>
            <ul className="list-disc pl-6 text-gray-700 space-y-2">
              <li>Use the Service for any illegal or unauthorized purpose</li>
              <li>Violate any laws in your jurisdiction</li>
              <li>Infringe upon the intellectual property rights of others</li>
              <li>Transmit any viruses, malware, or harmful code</li>
              <li>Attempt to gain unauthorized access to our systems</li>
              <li>Interfere with or disrupt the Service or servers</li>
              <li>Share or distribute content that is offensive, abusive, or inappropriate</li>
              <li>Use automated means to access the Service without permission</li>
            </ul>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">
              5. Intellectual Property Rights
            </h2>
            <p className="text-gray-700 mb-4">
              All content on EgalDeutsch, including stories, quizzes, text,
              graphics, logos, and software, is the property of EgalDeutsch or
              its content suppliers and is protected by international copyright,
              trademark, and other intellectual property laws.
            </p>
            <p className="text-gray-700">
              You may not copy, reproduce, distribute, modify, or create
              derivative works of any content from our Service without explicit
              written permission.
            </p>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">
              6. Third-Party Services
            </h2>
            <p className="text-gray-700 mb-4">
              Our Service uses third-party services to provide functionality and
              improve user experience:
            </p>
            <ul className="list-disc pl-6 text-gray-700 space-y-2">
              <li>
                <strong>Netlify:</strong> Our platform is hosted on Netlify's
                infrastructure for reliable and fast content delivery
              </li>
              <li>
                <strong>Google Services:</strong> We may use Google Analytics or
                other Google services for analytics and functionality
              </li>
              <li>
                <strong>MongoDB:</strong> For secure data storage and management
              </li>
            </ul>
            <p className="text-gray-700 mt-4">
              These third-party services have their own terms and privacy
              policies, which we encourage you to review.
            </p>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">
              7. GDPR Compliance and EU Users
            </h2>
            <p className="text-gray-700 mb-4">
              We comply with the General Data Protection Regulation (GDPR) and
              German data protection laws. As a user, you have the right to:
            </p>
            <ul className="list-disc pl-6 text-gray-700 space-y-2">
              <li>Access your personal data</li>
              <li>Rectify inaccurate personal data</li>
              <li>Request erasure of your personal data ("right to be forgotten")</li>
              <li>Restrict or object to processing of your personal data</li>
              <li>Data portability</li>
              <li>Withdraw consent at any time</li>
            </ul>
            <p className="text-gray-700 mt-4">
              For more information about data processing, please see our{" "}
              <Link to="/privacy-policy" className="text-blue-600 hover:underline">
                Privacy Policy
              </Link>
              .
            </p>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">
              8. Limitation of Liability
            </h2>
            <p className="text-gray-700 mb-4">
              To the maximum extent permitted by law, EgalDeutsch shall not be
              liable for any indirect, incidental, special, consequential, or
              punitive damages, or any loss of profits or revenues, whether
              incurred directly or indirectly, or any loss of data, use,
              goodwill, or other intangible losses.
            </p>
            <p className="text-gray-700">
              The Service is provided "as is" without warranties of any kind,
              either express or implied.
            </p>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">
              9. Termination
            </h2>
            <p className="text-gray-700">
              We may terminate or suspend your account and access to the Service
              immediately, without prior notice or liability, for any reason,
              including breach of these Terms. Upon termination, your right to
              use the Service will immediately cease.
            </p>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">
              10. Changes to Terms
            </h2>
            <p className="text-gray-700">
              We reserve the right to modify these Terms at any time. We will
              notify users of any material changes by posting the new Terms on
              this page and updating the "Last Updated" date. Your continued use
              of the Service after changes constitutes acceptance of the new
              Terms.
            </p>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">
              11. Governing Law
            </h2>
            <p className="text-gray-700">
              These Terms shall be governed by and construed in accordance with
              the laws of Germany and the European Union, without regard to its
              conflict of law provisions.
            </p>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">
              12. Contact Information
            </h2>
            <p className="text-gray-700 mb-4">
              If you have any questions about these Terms of Service, please
              contact us:
            </p>
            <div className="bg-gray-50 rounded-lg p-6">
              <p className="text-gray-700">
                <strong>Email:</strong>{" "}
                <a
                  href="mailto:egaldeutsch.com@gmail.com"
                  className="text-blue-600 hover:underline"
                >
                  egaldeutsch.com@gmail.com
                </a>
              </p>
              <p className="text-gray-700">
                <strong>Contact:</strong> Steve Phan
              </p>
            </div>
          </section>
        </div>

        <div className="mt-12 pt-8 border-t border-gray-200">
          <div className="text-center">
            <Link
              to="/"
              className="text-blue-600 hover:underline font-medium"
            >
              ← Back to Home
            </Link>
          </div>
        </div>
      </div>
    </Layout>
  );
};

export default TermsOfServicePage;
