
import Foundation

class APIService {
    static let shared = APIService()

    private let baseURL = URL(string: "http://localhost:8080")!

    func login(user: User, completion: @escaping (Result<Void, Error>) -> Void) {
        let url = baseURL.appendingPathComponent("login")
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.addValue("application/json", forHTTPHeaderField: "Content-Type")

        do {
            let jsonData = try JSONEncoder().encode(user)
            request.httpBody = jsonData
        } catch {
            completion(.failure(error))
            return
        }

        URLSession.shared.dataTask(with: request) { data, response, error in
            if let error = error {
                completion(.failure(error))
                return
            }

            guard let httpResponse = response as? HTTPURLResponse, (200...299).contains(httpResponse.statusCode) else {
                completion(.failure(NSError(domain: "Invalid response", code: 0, userInfo: nil)))
                return
            }

            completion(.success(()))
        }.resume()
    }

    func register(user: User, completion: @escaping (Result<Void, Error>) -> Void) {
        let url = baseURL.appendingPathComponent("register")
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.addValue("application/json", forHTTPHeaderField: "Content-Type")

        do {
            let jsonData = try JSONEncoder().encode(user)
            request.httpBody = jsonData
        } catch {
            completion(.failure(error))
            return
        }

        URLSession.shared.dataTask(with: request) { data, response, error in
            if let error = error {
                completion(.failure(error))
                return
            }

            guard let httpResponse = response as? HTTPURLResponse, (200...299).contains(httpResponse.statusCode) else {
                completion(.failure(NSError(domain: "Invalid response", code: 0, userInfo: nil)))
                return
            }

            completion(.success(()))
        }.resume()
    }

    func getConversations(completion: @escaping (Result<[Conversation], Error>) -> Void) {
        let url = baseURL.appendingPathComponent("conversation")
        var request = URLRequest(url: url)
        request.httpMethod = "GET"

        URLSession.shared.dataTask(with: request) { data, response, error in
            if let error = error {
                completion(.failure(error))
                return
            }

            guard let httpResponse = response as? HTTPURLResponse, (200...299).contains(httpResponse.statusCode) else {
                completion(.failure(NSError(domain: "Invalid response", code: 0, userInfo: nil)))
                return
            }

            guard let data = data else {
                completion(.failure(NSError(domain: "No data", code: 0, userInfo: nil)))
                return
            }

            do {
                let conversations = try JSONDecoder().decode([Conversation].self, from: data)
                completion(.success(conversations))
            } catch {
                completion(.failure(error))
            }
        }.resume()
    }
    
    func getMessages(conversationID: Int, completion: @escaping (Result<[Message], Error>) -> Void) {
        let url = baseURL.appendingPathComponent("conversation/\(conversationID)/message")
        URLSession.shared.dataTask(with: url) { data, response, error in
            if let error = error {
                completion(.failure(error))
                return
            }

            guard let data = data else {
                completion(.failure(NSError(domain: "No data", code: 0, userInfo: nil)))
                return
            }

            do {
                let messages = try JSONDecoder().decode([Message].self, from: data)
                completion(.success(messages))
            } catch {
                completion(.failure(error))
            }
        }.resume()
    }

    func sendMessage(conversationID: Int, content: String, completion: @escaping (Result<Void, Error>) -> Void) {
        let url = baseURL.appendingPathComponent("conversation/\(conversationID)/message")
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.addValue("application/json", forHTTPHeaderField: "Content-Type")

        let message = ["content": content]
        do {
            let jsonData = try JSONEncoder().encode(message)
            request.httpBody = jsonData
        } catch {
            completion(.failure(error))
            return
        }

        URLSession.shared.dataTask(with: request) { data, response, error in
            if let error = error {
                completion(.failure(error))
                return
            }

            guard let httpResponse = response as? HTTPURLResponse, (200...299).contains(httpResponse.statusCode) else {
                completion(.failure(NSError(domain: "Invalid response", code: 0, userInfo: nil)))
                return
            }

            completion(.success(()))
        }.resume()
    }
    
    func searchUsers(query: String, completion: @escaping (Result<[User], Error>) -> Void) {
        let url = baseURL.appendingPathComponent("user/search/\(query)")
        URLSession.shared.dataTask(with: url) { data, response, error in
            if let error = error {
                completion(.failure(error))
                return
            }

            guard let data = data else {
                completion(.failure(NSError(domain: "No data", code: 0, userInfo: nil)))
                return
            }

            do {
                let users = try JSONDecoder().decode([User].self, from: data)
                completion(.success(users))
            } catch {
                completion(.failure(error))
            }
        }.resume()
    }
    
    func createConversation(with userID: Int, completion: @escaping (Result<Conversation, Error>) -> Void) {
        let url = baseURL.appendingPathComponent("conversation")
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.addValue("application/json", forHTTPHeaderField: "Content-Type")

        let body = ["user_id": userID]
        do {
            let jsonData = try JSONEncoder().encode(body)
            request.httpBody = jsonData
        } catch {
            completion(.failure(error))
            return
        }

        URLSession.shared.dataTask(with: request) { data, response, error in
            if let error = error {
                completion(.failure(error))
                return
            }

            guard let data = data else {
                completion(.failure(NSError(domain: "No data", code: 0, userInfo: nil)))
                return
            }

            do {
                let conversation = try JSONDecoder().decode(Conversation.self, from: data)
                completion(.success(conversation))
            } catch {
                completion(.failure(error))
            }
        }.resume()
    }
}
