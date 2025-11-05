
import SwiftUI

struct LoginView: View {
    @Binding var isAuthenticated: Bool
    @State private var username = ""
    @State private var password = ""
    @State private var isRegistering = false

    var body: some View {
        VStack {
            Text(isRegistering ? "Register" : "Login")
                .font(.largeTitle)
                .padding()

            TextField("Username", text: $username)
                .textFieldStyle(RoundedBorderTextFieldStyle())
                .padding()

            SecureField("Password", text: $password)
                .textFieldStyle(RoundedBorderTextFieldStyle())
                .padding()

            Button(action: {
                let user = User(Username: username, Password: password)
                if isRegistering {
                    APIService.shared.register(user: user) { result in
                        switch result {
                        case .success:
                            // For simplicity, automatically log in after registration
                            APIService.shared.login(user: user) { result in
                                switch result {
                                case .success:
                                    isAuthenticated = true
                                case .failure(let error):
                                    print(error.localizedDescription)
                                }
                            }
                        case .failure(let error):
                            print(error.localizedDescription)
                        }
                    }
                } else {
                    APIService.shared.login(user: user) { result in
                        switch result {
                        case .success:
                            isAuthenticated = true
                        case .failure(let error):
                            print(error.localizedDescription)
                        }
                    }
                }
            }) {
                Text(isRegistering ? "Register" : "Login")
            }
            .padding()

            Button(action: {
                isRegistering.toggle()
            }) {
                Text(isRegistering ? "Already have an account? Login" : "Don't have an account? Register")
            }
            .padding()
        }
    }
}

struct LoginView_Previews: PreviewProvider {
    static var previews: some View {
        LoginView(isAuthenticated: .constant(false))
    }
}
